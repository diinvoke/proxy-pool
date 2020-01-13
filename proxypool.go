package proxypool

import (
	"github.com/mingcheng/proxypool/controller"
	"github.com/mingcheng/proxypool/model"
	"github.com/mingcheng/proxypool/spider"
	"github.com/mingcheng/proxypool/storage"
	"math/rand"
	"time"
)

var (
	spiderTicker *time.Ticker
	proxyChecker *controller.Checker
)

// Config type for config the proxypool
type Config struct {
	FetchInterval   time.Duration
	CheckInterval   time.Duration
	CheckConcurrent uint
}

// Start proxypool
func Start(config Config) {
	proxyChecker = &controller.Checker{
		Storage: storage.NewMemoryStorage(),
	}

	proxyChecker.Start(config.CheckInterval, config.CheckConcurrent)
	spiderTicker = spider.InitSpider(config.FetchInterval, *proxyChecker)
}

// Stop proxypool
func Stop() {
	proxyChecker.Stop()
	spiderTicker.Stop()
}

// Random to get a validable proxy
func Random() *model.Proxy {
	// @see https://stackoverflow.com/questions/33994677/pick-a-random-value-from-a-go-slice
	if all := All(); len(all) > 0 {
		rand.Seed(time.Now().Unix())
		return all[rand.Intn(len(all))]
	}

	return nil
}

// All function to get all valid proxy
func All() []*model.Proxy {
	return proxyChecker.Storage.All()
}

// Add a proxy to check
func Add(p model.Proxy) {
	go proxyChecker.Check(&p)
}
