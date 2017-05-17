package util

import "testing"

func TestFloatEqual(t *testing.T) {
	if FloatEqual(-0.005422689076629065, 0) == false {
		t.Error("")
	}
}
