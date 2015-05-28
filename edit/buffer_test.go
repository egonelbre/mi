package edit

import "testing"

func TestRegion(t *testing.T) {
	r := Region{}
	if !r.Overlaps(&r) {
		t.Error("region should overlap with itself")
	}
}
