package network

import "testing"

func TestIpPrintint(t *testing.T) {
	IP := IPAddr{8, 8, 8, 8}
	IPExpected := "8.8.8.8"
	IPReturned := IP.String()
	if IPExpected != IPReturned {
		t.Errorf("expected %s, got %s", IPExpected, IPReturned)
	}
}
