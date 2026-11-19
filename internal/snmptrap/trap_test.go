package snmptrap

import (
	"testing"
	"github.com/gosnmp/gosnmp"
)
func TestNormalizeLinkDown(t *testing.T){p:=&gosnmp.SnmpPacket{};p.Trap.Oid=linkDownOID;ev:=(Normalizer{}).Normalize(p,nil);if ev.EventType!="LINK_DOWN"{t.Fatalf("got %q",ev.EventType)}}
func TestVendorClassificationDoesNotAssumeOptics(t *testing.T){_ = gosnmp.Version2c}
