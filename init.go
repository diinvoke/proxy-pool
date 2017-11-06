package proxypool

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/Agzdjy/proxy-pool/spider"
	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/go-redis/redis"
)

const configFile = "./config.json"

var ip181Spider = &spider.Ip181{}

func readJson() map[string]string {
	jsonData := map[string]string{}
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic("config is undefined")
	}

	err = json.Unmarshal(bytes, &jsonData)
	if err != nil {
		panic("json format error")
	}
	return jsonData
}

func initStorage() storage.Storage {
	config := readJson()
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
	ip181.Do("http://www.ip181.com/", stor)
}
