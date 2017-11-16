package util

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

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
		Timeout:   time.Millisecond * 500,
	}

	resq, err := client.Get("http://ip.gs/")
	if err != nil {
		return false
	}
	defer resq.Body.Close()
	return true
}
