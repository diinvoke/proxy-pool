package storage

import (
	"testing"

	"github.com/Agzdjy/proxy-pool/model"
	"github.com/go-redis/redis"
)

var storage Storage = NewRedisClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})

var ip = &model.IP{
	Address:  "192.168.1.20",
	Port:     "6908",
	Protocol: "https",
}

func TestSave(t *testing.T) {
	err := storage.Save(ip)

	if err != nil {
		t.Error("save error", err)
	}
}

func TestRange(t *testing.T) {
	TestSave(t)

	rangeIp, err := storage.RangeOne("https")
	if err != nil {
		t.Error("range key error", err)
		return
	}

	if rangeIp.Address == "" {
		t.Error("range no one")
		return
	}
}

func TestDel(t *testing.T) {
	TestRange(t)

	ok := storage.Del("https")
	if !ok {
		t.Error("del failed")
	}
}
