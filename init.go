package proxypool

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"sync"

	"github.com/Agzdjy/proxy-pool/spider"
	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/go-redis/redis"
)

func readJson(configPath string) map[string]string {
	jsonData := map[string]string{}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic("config is undefined")
	}

	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		panic("json format error")
	}
	return jsonData
}

func initStorage(configPath string) storage.Storage {
	config := readJson(configPath)
	store := config["store"]
	var stor storage.Storage

	switch store {
	case "redis":
		db, _ := strconv.Atoi(config["db"])
		addr := config["address"]
		password := config["password"]
		stor = storage.NewRedisClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})
	}
	return stor
}

func initSpider(stor storage.Storage) {
	var ip181 spider.Spider = &spider.Ip181{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		ip181.Do("http://www.ip181.com/", stor)
		wg.Done()
	}()

	wg.Wait()
}
