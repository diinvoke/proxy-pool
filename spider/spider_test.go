package spider

import (
	"testing"

	"github.com/diinvoke/proxy-pool/model"
	"github.com/diinvoke/proxy-pool/storage"
)

func TestDo(t *testing.T) {
	storage := storage.NewLocalCache()
	ip181 := NewIP181(storage)
	err := ip181.Do()

	t.Logf("count:%d", ip181.LoadCount())

	if err != nil {
		t.Error("spider do error", err)
		return
	}
	ip, _ := storage.Random(model.ProtocolHttp)
	if ip == nil {
		t.Errorf("got nil")
		return
	}

	t.Logf("%s://%s:%s", ip.Protocol, ip.Address, ip.Port)
}
