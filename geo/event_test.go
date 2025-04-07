package geo

import "testing"

func TestEventXray(t *testing.T) {
	flareC := XRAY_FLARE_C
	if flareC.String() != "Class C" {
		t.Errorf("Expected 'Class C', got '%s'", flareC.String())
	}
}
