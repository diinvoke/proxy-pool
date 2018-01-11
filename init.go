package proxypool

import (
	"strconv"

	"sync"

	"github.com/Agzdjy/proxy-pool/spider"
	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/Agzdjy/proxy-pool/util"
	"github.com/go-redis/redis"
)

func getStorage(configPath string) storage.Storage {
	var store storage.Storage

	config := util.ReadJson(configPath)
	storeType := config["store"]

	switch storeType {
	case "redis":
		db, _ := strconv.Atoi(config["db"])
		addr := config["address"]
		password := config["password"]
		store = storage.NewRedisClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	default:
		panic("unknown store type")
	}
	return store
}

var spiderSource = map[string]spider.Spider{
	"http://www.ip181.com/": &spider.Ip181{},
}

func initSpider(store storage.Storage) {
	var wg sync.WaitGroup

	for url, spider := range spiderSource {
		wg.Add(1)
		go func() {
			spider.Do(url, store)
			wg.Done()
		}()

	}

	wg.Wait()
}
