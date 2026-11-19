package alerts

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/roky61980-create/port-monitor/internal/events"
	"github.com/roky61980-create/port-monitor/internal/model"
	"github.com/roky61980-create/port-monitor/internal/store"
)

type Engine struct { db *store.Store; hub *events.Hub }
func New(db *store.Store, hub *events.Hub) *Engine { return &Engine{db:db,hub:hub} }

func (e *Engine) Process(ev model.NetworkEvent) error {
	if ev.Timestamp.IsZero() { ev.Timestamp=time.Now().UTC() }
	if ev.ID=="" { ev.ID=id() }
	if err:=e.db.SaveEvent(ev); err!=nil { return err }
	if e.hub!=nil { e.hub.Publish(map[string]any{"kind":"event","event":ev}) }
	key:=strings.ToUpper(ev.EventType)
	switch key {
	case "LINK_DOWN","PORT_DOWN","DEVICE_DOWN","OPTICAL_WARNING","OPTICAL_CRITICAL":
		severity:=ev.Severity; if severity=="" { severity="CRITICAL" }; a,err:=e.db.AlertByKey(ev.DeviceID,ev.Interface,key)
		if err==nil { a.UpdatedAt=ev.Timestamp; a.CurrentValue=ev.Message; if severity=="CRITICAL" { a.Severity=severity }; if err=e.db.SaveAlert(a);err!=nil{return err} } else if errors.Is(err,store.ErrNotFound) {
			a=model.Alert{ID:id(),DeviceID:ev.DeviceID,Device:ev.Device,Interface:ev.Interface,AlertType:key,Severity:severity,Status:"OPEN",CurrentValue:ev.Message,CreatedAt:ev.Timestamp,UpdatedAt:ev.Timestamp}
			if err=e.db.SaveAlert(a);err!=nil{return err}; if e.hub!=nil { e.hub.Publish(map[string]any{"kind":"alert","event_type":"ALERT_CREATED","alert":a}) }
		} else { return err }
	case "LINK_UP","PORT_UP","DEVICE_UP","OPTICAL_RESOLVED":
		alerts,err:=e.db.Alerts(); if err!=nil{return err}
		for _,a:=range alerts { if a.DeviceID==ev.DeviceID && a.Interface==ev.Interface && a.Status!="RESOLVED" && related(a.AlertType,key) { now:=ev.Timestamp; a.Status="RESOLVED"; a.ResolvedAt=&now; a.UpdatedAt=now; if err=e.db.SaveAlert(a);err!=nil{return err}; if e.hub!=nil { e.hub.Publish(map[string]any{"kind":"alert","event_type":"ALERT_RESOLVED","alert":a}) } } }
	}
	return nil
}
func related(a,b string) bool { return (strings.Contains(a,"PORT")&&strings.Contains(b,"PORT"))||(strings.Contains(a,"LINK")&&strings.Contains(b,"LINK"))||(strings.Contains(a,"DEVICE")&&strings.Contains(b,"DEVICE"))||(strings.Contains(a,"OPTICAL")&&strings.Contains(b,"OPTICAL")) }
func id() string { var b [16]byte; if _,err:=rand.Read(b[:]);err!=nil{return time.Now().UTC().Format("20060102150405.000000000")};return hex.EncodeToString(b[:]) }
