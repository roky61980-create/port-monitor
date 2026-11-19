package snmptrap

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/roky61980-create/port-monitor/internal/model"
)

// Standard linkDown/linkUp notification OIDs from IF-MIB are used only for event classification.
const (
	linkDownOID = ".1.3.6.1.6.3.1.1.5.3"
	linkUpOID = ".1.3.6.1.6.3.1.1.5.4"
)

type Resolver interface { DeviceByIP(ip string) (model.Device, error); PortByIndex(deviceID string, index int) (model.Port, error) }

type Normalizer struct { Devices Resolver }

func (n Normalizer) Normalize(packet *gosnmp.SnmpPacket, addr *net.UDPAddr) model.NetworkEvent {
	now:=time.Now().UTC(); ev:=model.NetworkEvent{Timestamp:now,Severity:"INFO",EventType:"DEVICE_EVENT",Message:"SNMP notification received",OID:"",RawData:map[string]any{"source":addr.String()}}
	if packet==nil{return ev}
	vars:=make([]map[string]any,0,len(packet.Variables)); var ifIndex int
	trapOID:=packet.Trap.Oid
	for _,v:=range packet.Variables { oid:=strings.TrimPrefix(v.Name,"."); value:=fmt.Sprint(v.Value); vars=append(vars,map[string]any{"oid":oid,"type":v.Type.String(),"value":value}); if strings.HasSuffix(oid,".1.3.6.1.2.1.2.2.1.1") { if x,ok:=asInt(v.Value);ok {ifIndex=x} }; if oid=="1.3.6.1.6.3.1.1.4.1.0" { trapOID=fmt.Sprint(v.Value) } }
	ev.OID=trapOID; ev.RawData=map[string]any{"source":addr.String(),"varbinds":vars}
	switch strings.TrimPrefix(trapOID,".") { case strings.TrimPrefix(linkDownOID,"."): ev.EventType="LINK_DOWN";ev.Severity="CRITICAL";ev.Message="Interface link down"; case strings.TrimPrefix(linkUpOID,"."): ev.EventType="LINK_UP";ev.Message="Interface link up" }
	if n.Devices!=nil { if d,err:=n.Devices.DeviceByIP(addr.IP.String());err==nil {ev.DeviceID=d.ID;ev.Device=d.Name;ev.Vendor=d.Vendor;ev.Model=d.Model} }
	ev.InterfaceID=ifIndex
	if ifIndex>0&&n.Devices!=nil&&ev.DeviceID!="" { if p,err:=n.Devices.PortByIndex(ev.DeviceID,ifIndex);err==nil {ev.Interface=p.IfName;if ev.Interface==""{ev.Interface=p.IfDescr}} }
	return ev
}
func asInt(v any)(int,bool){switch x:=v.(type){case int:return x,true;case int32:return int(x),true;case uint32:return int(x),true;case int64:return int(x),true;case uint64:return int(x),true;case []byte:var n int;_,err:=fmt.Sscanf(string(x),"%d",&n);return n,err==nil;default:var n int;_,err:=fmt.Sscanf(fmt.Sprint(v),"%d",&n);return n,err==nil}}

type Receiver struct { listener *gosnmp.TrapListener }
func Start(listen string, normalizer Normalizer, emit func(model.NetworkEvent) error) (*Receiver,error) {
	l:=gosnmp.NewTrapListener(); l.Params=&gosnmp.GoSNMP{Version:gosnmp.Version2c,Community:community(),Timeout:3*time.Second,Retries:1}
	l.OnNewTrap=func(p *gosnmp.SnmpPacket,a *net.UDPAddr){ ev:=normalizer.Normalize(p,a); if err:=emit(ev);err!=nil { fmt.Printf("trap processing failed: %v\n",err) } }
	if err:=l.Listen(listen);err!=nil{return nil,err};return &Receiver{listener:l},nil
}
func (r *Receiver) Close() error { if r==nil||r.listener==nil{return nil};return r.listener.Close() }
func community() string { if v:=getenv("SNMP_COMMUNITY");v!=""{return v};return "public" }
