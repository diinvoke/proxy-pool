package model

type Protocol string

const (
	ProtocolHttps = "https"
	ProtocolHttp  = "http"
)

func ProtocolFromString(s string) Protocol {
	switch s {
	case ProtocolHttp, ProtocolHttps:
		return Protocol(s)
	}

	return Protocol("")
}

func ProtocolToString(protocol Protocol) string {
	return string(protocol)
}

type IP struct {
	Address  string   `json:"address"`
	Port     string   `json:"port"`
	Protocol Protocol `json:"protocol"`
}

func (i *IP) Equal(ip *IP) bool {
	if ip == nil || i == nil {
		return false
	}

	return i.Address == ip.Address &&
		i.Protocol == ip.Protocol &&
		i.Port == ip.Port
}
