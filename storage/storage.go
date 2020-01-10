package storage

import (
	"github.com/mingcheng/proxypool/model"
	rpc "github.com/mingcheng/proxypool/protobuf"
)

type IStorage interface {
	Save(proxy *model.Proxy) error
	Del(proxy *model.Proxy) bool
	Random(protocol rpc.Protocol) (*model.Proxy, error)
	All() []*model.Proxy
	Close() error
}
