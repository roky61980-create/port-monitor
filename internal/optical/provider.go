package optical

import (
	"context"
	"sync"
	"time"

	"github.com/roky61980-create/port-monitor/internal/model"
)

type FieldStatus string
const ( Verified FieldStatus="VERIFIED"; NotVerified FieldStatus="NOT_VERIFIED"; Unavailable FieldStatus="UNAVAILABLE" )
type Sensor struct { Interface string; Metric string; OID string; Status FieldStatus; Scale float64; Unit string }
type Provider interface { Discover(context.Context,model.Device)([]Sensor,error); Read(context.Context,model.Device,[]Sensor)([]model.OpticalReading,error) }
type Registry struct { providers map[string]Provider }
func NewRegistry() *Registry { return &Registry{providers:map[string]Provider{"cisco":UnsupportedProvider{},"mikrotik":UnsupportedProvider{},"bdcom":UnsupportedProvider{},"generic":UnsupportedProvider{},"unknown":UnsupportedProvider{}}} }
func (r *Registry) For(vendor string) Provider { if p:=r.providers[normalize(vendor)];p!=nil{return p};return r.providers["unknown"] }
func normalize(v string)string { switch {case contains(v,"cisco"):return "cisco";case contains(v,"mikrotik"):return "mikrotik";case contains(v,"bdcom"):return "bdcom";case v=="":return "unknown";default:return "generic"} }
func contains(a,b string)bool { for i:=0;i+len(b)<=len(a);i++{if a[i:i+len(b)]==b{return true}};return false }

type UnsupportedProvider struct{}
func (UnsupportedProvider) Discover(context.Context,model.Device)([]Sensor,error) { return []Sensor{{Metric:"rx_power",Status:NotVerified},{Metric:"tx_power",Status:NotVerified},{Metric:"temperature",Status:NotVerified},{Metric:"voltage",Status:NotVerified},{Metric:"bias_current",Status:NotVerified}},nil }
func (UnsupportedProvider) Read(context.Context,model.Device,[]Sensor)([]model.OpticalReading,error) { return nil,nil }

// PollDevices bounds concurrency and does not overlap polls for the same device.
func PollDevices(ctx context.Context, devices []model.Device, workers int, interval time.Duration, p func(context.Context,model.Device)error) {
	if workers<1 {workers=1}; sem:=make(chan struct{},workers); var wg sync.WaitGroup; ticker:=time.NewTicker(interval);defer ticker.Stop()
	for { select {case <-ctx.Done():wg.Wait();return;case <-ticker.C:
		for _,d:=range devices { select {case sem<-struct{}{}:wg.Add(1);go func(dev model.Device){defer wg.Done();defer func(){<-sem}();_ = p(ctx,dev)}(d);case <-ctx.Done():wg.Wait();return} }
	} }
}
