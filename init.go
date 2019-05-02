package proxypool

import (
	"fmt"
	"sync"
	"time"

	"github.com/diinvoke/proxy-pool/spider"
	"github.com/diinvoke/proxy-pool/storage"
)

func initProxyPool(store storage.IStorage) {
	spiders := []spider.ISpider{
		spider.NewIP181(store),
	}

	var wg sync.WaitGroup
	wg.Add(len(spiders))

	for _, sp := range spiders {
		go func(group *sync.WaitGroup, iSpider spider.ISpider) {
			err := iSpider.Do()
			if err != nil {
				fmt.Println("name, load failed, err", iSpider.Name(), err)
			}
			fmt.Println("name:, load success:", iSpider.Name(), iSpider.LoadCount())

			group.Done()
		}(&wg, sp)
	}

	wg.Wait()
}

func autoLoad() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			initProxyPool(store)
		}
	}
}
