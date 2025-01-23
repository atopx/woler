package woler

import "testing"

func TestDo(t *testing.T) {
	var macaddr = "F0-2F-74-B0-1D-E0"
	if err := Do(macaddr); err != nil {
		t.Error(err)
	}
}

func TestDoWithBroadcast(t *testing.T) {
	var macaddr = "F0-2F-74-B0-1D-E0"
	if err := DoWithBroadcast(macaddr, "255.255.255.255", 9); err != nil {
		t.Error(err)
	}
}
