package storage

import (
	"github.com/diinvoke/proxy-pool/model"
)

type IStorage interface {
	Save(ip *model.IP) error
	Del(ip *model.IP) bool
	Random(protocol model.Protocol) (*model.IP, error)
	Close() error
}
