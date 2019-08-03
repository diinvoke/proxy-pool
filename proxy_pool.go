package proxypool

import (
	"github.com/diinvoke/proxy-pool/model"
	"github.com/diinvoke/proxy-pool/storage"
	"github.com/diinvoke/proxy-pool/util"
)

var store storage.IStorage

func init() {
	store = storage.NewLocalCache()
	initProxyPool(store)

	go autoLoad()
}

// Range one proxy ip
func Range(protocol string) *model.IP {
	ip, err := store.Random(model.ProtocolFromString(protocol))
	if err != nil {
		return nil
	}
	return ip
}

// Check proxy valid
func Check(proxy string) bool {
	return util.Check(proxy)
}

// Del failure proxy
func Del(ip *model.IP) bool {
	return store.Del(ip)
}

func RefreshData() {
	initProxyPool(store)
}
