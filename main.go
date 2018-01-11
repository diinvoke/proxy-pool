package proxypool

import (
	"time"

	"github.com/Agzdjy/proxy-pool/model"
	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/Agzdjy/proxy-pool/util"
)

var store storage.Storage = nil

// init storage and spider
func InitData(configPath string) {
	store = getStorage(configPath)
	initSpider(store)

	ticker := time.NewTicker(3 * time.Minute)
	go func() {
		for _ = range ticker.C {
			initSpider(store)
		}
	}()
}

// Range one proxy ip
func Range(protocol string) *model.IP {
	ip, err := store.RangeOne(protocol)
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
	initSpider(store)
}
