package storage

import (
	"github.com/mingcheng/proxypool/model"
)

type IStorage interface {
	Save(proxy *model.Proxy) error
	Del(proxy *model.Proxy) bool
	Random(protocol model.Protocol) (*model.Proxy, error)
	All() []*model.Proxy
	Close() error
}
