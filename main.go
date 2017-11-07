package proxypool

import (
	"time"

	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/Agzdjy/proxy-pool/model"
	"github.com/Agzdjy/proxy-pool/storage"
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
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		panic("proxy schema error")
	}

	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Millisecond * 200,
	}

	resq, err := client.Get("http://ip.gs/")
	if err != nil {
		return false
	}
	defer resq.Body.Close()
	return true
}

// Del failure proxy
func Del(ip *model.IP) bool {
	return stor.Del(ip)
}
