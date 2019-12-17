package model

import (
	"crypto/tls"
	"fmt"
	simpleJSON "github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Protocol string

const (
	ProtocolHttps = "https"
	ProtocolHttp  = "http"
	ProtocolSocks = "socks"
)

func ProtocolFromString(s string) Protocol {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case ProtocolHttp, ProtocolHttps, ProtocolSocks:
		return Protocol(s)
	}

	return Protocol("")
}

func ProtocolToString(protocol Protocol) string {
	return string(protocol)
}

type Proxy struct {
	Address   string    `json:"address"`
	Port      int       `json:"port"`
	Protocol  Protocol  `json:"protocol"`
	LastCheck time.Time `json:"last_check"`
	Speed     int64     `json:"speed"`
	From      string    `json:"from"`
}

func (i *Proxy) Equal(ip *Proxy) bool {
	if ip == nil || i == nil {
		return false
	}
	return i.Address == ip.Address && i.Protocol == ip.Protocol && i.Port == ip.Port
}

func (i *Proxy) String() string {
	return fmt.Sprintf("%s://%s:%d", i.Protocol, i.Address, i.Port)
}

func (i *Proxy) Check() (bool, error) {
	log.Printf("start check %s", i.String())
	begin := time.Now()

	i.LastCheck = time.Now()
	proxyUrl, err := url.Parse(i.String())
	if err != nil {
		return false, fmt.Errorf("invalid schema format, %s", err)
	}

	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxyUrl),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("proxy %s check failed, %s", i.String(), err)
	}
	defer resp.Body.Close()

	// check return data
	if _, err := simpleJSON.NewFromReader(resp.Body); err != nil {
		return false, err
	}

	// mark speeds
	i.Speed = time.Now().Sub(begin).Nanoseconds()
	log.Printf("check finished, %s is valid proxy", i.String())
	return true, nil
}
