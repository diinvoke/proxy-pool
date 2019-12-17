package storage

import (
	"github.com/mingcheng/proxypool/model"
	"github.com/patrickmn/go-cache"
	"sync"
)

import "time"

var (
	singleton *MemoryStorage
	once      sync.Once
)

type MemoryStorage struct {
	cache *cache.Cache
}

func NewMemoryStorage() *MemoryStorage {
	once.Do(func() {
		singleton = &MemoryStorage{
			cache: cache.New(30*time.Minute, 24*time.Hour),
		}
	})

	return singleton
}

func (m *MemoryStorage) Save(proxy *model.Proxy) error {
	m.cache.Set(proxy.String(), proxy, cache.DefaultExpiration)
	return nil
}

func (m *MemoryStorage) Del(proxy *model.Proxy) bool {
	m.cache.Delete(proxy.String())
	return true
}

func (*MemoryStorage) Random(protocol model.Protocol) (*model.Proxy, error) {
	return nil, nil
}

func (m *MemoryStorage) All() []*model.Proxy {
	var result []*model.Proxy
	for _, v := range m.cache.Items() {
		result = append(result, (v.Object).(*model.Proxy))
	}

	return result
}

func (m *MemoryStorage) Close() error {
	m.cache.Flush()
	return nil
}
