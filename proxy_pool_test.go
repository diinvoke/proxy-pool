package proxypool

import "testing"

func TestRange(t *testing.T) {
	rangeIp := Range("http")
	if rangeIp == nil || rangeIp.Address == "" {
		t.Error("range no one")
		return
	}
}
