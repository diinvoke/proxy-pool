package proxypool

import (
	"github.com/Agzdjy/proxy-pool/model"
	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/Agzdjy/proxy-pool/util"
)

var stor storage.Storage = nil

func init() {
	stor = initStorage()
	initSpider(stor)
}

// Range one proxy ip
func Range(protocol string) *model.IP {
	ip, err := stor.RangeOne(protocol)
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
	return stor.Del(ip)
}
