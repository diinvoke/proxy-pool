package storage

import (
	"github.com/Agzdjy/proxy-pool/model"
	"github.com/go-redis/redis"
)

type Storage interface {
	Save(ip *model.IP) error
	Del(ip *model.IP) bool
	RangeOne(protocol string) (ip *model.IP, err error)
	Close() error
}

func NewRedisClient(options *redis.Options) *Redis {
	return &Redis{client: redis.NewClient(options)}
}
