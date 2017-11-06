package proxypool

import "testing"

func TestRange(t *testing.T) {
	rangeIp := Range("https")

	if rangeIp == nil || rangeIp.Address == "" {
		t.Error("range no one")
		return
	}
}

func TestCheck(t *testing.T) {
	failed := Check("192.168.1.9:400")

	rangeIp := Range("http")
	success := Check(rangeIp.Protocol + "://" + rangeIp.Address + ":" + rangeIp.Port)

	if failed {
		t.Error("check ip failed")
	}
	if !success {
		t.Errorf("check ip failed ip: %s", rangeIp.Address)
	}

}
