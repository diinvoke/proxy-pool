package spider

import (
	"fmt"
	"sync"
	"time"

	"github.com/mingcheng/proxypool/controller"
)

func InitSpider(d time.Duration, proxyChecker controller.Checker) *time.Ticker {
	// @TODO
	spiders := []Spider{
		&Feiyi{},
		&IP89{},
	}

	ticker := time.NewTicker(d)
	for ; true; <-ticker.C {
		go func() {
			var wg sync.WaitGroup
			wg.Add(len(spiders))

			for _, sp := range spiders {
				go func(group *sync.WaitGroup, spider Spider) {
					proxies, err := spider.Do()
					if err != nil {
						fmt.Println("name, load failed, err", spider.Name(), err)
					}

					for _, v := range proxies {
						go proxyChecker.Check(v)
					}

					group.Done()
				}(&wg, sp)
			}

			wg.Wait()
		}()
	}

	return ticker
}
