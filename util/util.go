package util

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
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
		Timeout:   time.Millisecond * 500, //超时时间
	}

	resp, err := client.Get("http://ip.sb/")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	defer resp.Body.Close()
	return true
}

func HttpGet(url string) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Minute * 1,
	}
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

func ReadJson(jsonPath string) map[string]string {
	jsonMap := map[string]string{}
	bytes, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		panic("config is undefined")
	}

	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		panic("json format error")
	}
	return jsonMap
}
