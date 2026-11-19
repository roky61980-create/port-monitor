package optical

import "testing"
func TestUnsupportedValuesAreNotMeasurements(t *testing.T){p:=UnsupportedProvider{};s,e:=p.Discover(t.Context(),structDevice());if e!=nil{t.Fatal(e)};for _,v:=range s{if v.Status!=NotVerified{t.Fatalf("%s status %s",v.Metric,v.Status)}}}
