package storage

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"

	"github.com/diinvoke/proxy-pool/model"
)

func getProtocol() string {
	protocol := "http"
	if rand.Intn(10) < 5 {
		protocol = "https"
	}
	return protocol
}

func genIPModel() *model.IP {
	port := fmt.Sprintf("%d", rand.Int63n(1000))

	return &model.IP{
		Address:  "0.0.0.0",
		Port:     port,
		Protocol: model.Protocol(getProtocol()),
	}
}

func save(cache IStorage, count int) {
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(wgg *sync.WaitGroup) {
			ipModel := genIPModel()
			cache.Save(ipModel)
			wgg.Done()
		}(&wg)
	}

	wg.Wait()

	fmt.Println("save success")
}

func randomAndDel(cache IStorage, count int, t *testing.T) {
	cnt := count / 4
	var wg sync.WaitGroup
	wg.Add(cnt)

	for i := 0; i < cnt; i++ {
		go func(index int, wgg *sync.WaitGroup) {
			defer wgg.Done()

			ip, _ := cache.Random(model.Protocol(getProtocol()))
			if ip == nil {
				t.Errorf("random result is nil")
			}

			cache.Del(ip)
			lc := cache.(*LocalCache)
			if lc.Exist(ip) {
				t.Errorf("del result failed")
			}

			//fmt.Printf("%s://%s:%s\n", ip.Protocol, ip.Address, ip.Port)

		}(i, &wg)
	}

	wg.Wait()

	fmt.Println("randomAndDel success")
}

func TestLogic(t *testing.T) {
	cache := NewLocalCache()
	count := 2000

	save(cache, count)

	randomAndDel(cache, count, t)
}
