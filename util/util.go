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

func HttpGet(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
