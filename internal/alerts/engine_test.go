package alerts

import "testing"
func TestRelatedCorrelation(t *testing.T){if !related("PORT_DOWN","PORT_UP"){t.Fatal("port transitions should correlate")};if related("PORT_DOWN","DEVICE_UP"){t.Fatal("unrelated classes must not correlate")}}
