package model

import (
	"crypto/tls"
	"fmt"
	simpleJSON "github.com/bitly/go-simplejson"
	"github.com/golang/protobuf/ptypes"
	rpc "github.com/mingcheng/proxypool/protobuf"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//type Protocol string

const (
	ProtocolHttps   = "https"
	ProtocolHttp    = "http"
	ProtocolSocks   = "socks"
	ProtocolUnknown = "unknown"
)

func ProtocolFromString(s string) rpc.Protocol {
	s = strings.TrimSpace(strings.ToLower(s))
	switch s {
	case ProtocolHttp:
		return rpc.Protocol_HTTP
	case ProtocolHttps:
		return rpc.Protocol_HTTPS
	case ProtocolSocks:
		return rpc.Protocol_SOCKS
	}

	return rpc.Protocol_UNKNOWN
}

func ProtocolToString(protocol rpc.Protocol) string {
	switch protocol {
	case rpc.Protocol_HTTPS:
		return ProtocolHttps
	case rpc.Protocol_HTTP:
		return ProtocolHttp
	case rpc.Protocol_SOCKS:
		return ProtocolSocks
	}
	return ProtocolUnknown
}

type Proxy struct {
	rpc.Proxy
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

	i.LastCheck, _ = ptypes.TimestampProto(time.Now())
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
	i.Speed = uint64(time.Now().Sub(begin).Nanoseconds())
	log.Printf("check finished, %s is valid proxy", i.String())
	return true, nil
}
