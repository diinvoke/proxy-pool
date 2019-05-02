package storage

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/diinvoke/proxy-pool/model"
)

const (
	keyTMPl = "ip_proxy:%s"
)

type LocalCache struct {
	items map[string][]*model.IP
	sync.RWMutex
}

var _ IStorage = new(LocalCache)

func NewLocalCache() IStorage {
	return &LocalCache{
		items: make(map[string][]*model.IP, 2),
	}
}

func (l *LocalCache) makeKey(protocol model.Protocol) string {
	return fmt.Sprintf(keyTMPl, protocol)
}

func (l *LocalCache) MSave(ips []*model.IP) error {
	l.Lock()
	for _, ip := range ips {
		if ip != nil {
			key := l.makeKey(ip.Protocol)
			l.items[key] = append(l.items[key], ip)
		}
	}
	l.Unlock()

	return nil
}

func (l *LocalCache) Save(ip *model.IP) error {
	if ip == nil {
		return errors.New("ip is nil model")
	}

	return l.MSave([]*model.IP{ip})
}

func (l *LocalCache) Del(ip *model.IP) bool {
	if ip == nil {
		return true
	}

	l.Lock()

	ips := l.items[l.makeKey(ip.Protocol)]
	pruneIps := make([]*model.IP, 0, len(ips)-1)

	for _, i := range ips {
		if !i.Equal(ip) {
			pruneIps = append(pruneIps, i)
		}
	}

	l.items[l.makeKey(ip.Protocol)] = pruneIps

	l.Unlock()

	return true
}

func (l *LocalCache) Random(protocol model.Protocol) (*model.IP, error) {
	l.RLock()
	defer l.RUnlock()

	ips := l.items[l.makeKey(protocol)]

	if len(ips) == 0 {
		return nil, fmt.Errorf("empty")
	}
	index := rand.Intn(len(ips))

	return ips[index], nil
}

func (l *LocalCache) Clear() {
	l.Lock()
	l.items = make(map[string][]*model.IP, 2)
	l.Unlock()
}

func (l *LocalCache) Exist(ip *model.IP) bool {
	if ip == nil {
		return false
	}

	l.RLock()
	defer l.RUnlock()

	ips := l.items[l.makeKey(ip.Protocol)]
	for _, i := range ips {
		if !i.Equal(ip) {
			continue
		}

		return true
	}

	return false
}

func (l *LocalCache) Close() error {
	return nil
}
