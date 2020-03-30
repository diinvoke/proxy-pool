package controller

import (
	"fmt"
	"github.com/mingcheng/proxypool/model"
	"github.com/mingcheng/proxypool/storage"

	"log"
	"time"
)

type Checker struct {
	checkTicker *time.Ticker
	checkChan   chan *model.Proxy
	done        chan bool
	Storage     storage.IStorage
}

func (c *Checker) Check(proxy *model.Proxy) {
	if c.checkChan != nil {
		c.checkChan <- proxy
	}
}

func (c *Checker) Start(interval time.Duration, concurrent uint) {
	c.checkChan = make(chan *model.Proxy, concurrent)
	c.checkTicker = time.NewTicker(interval)
	c.done = make(chan bool)

	go func() {
		for {
			select {
			case <-c.done:
				return
			case proxy := <-c.checkChan:
				if result, err := proxy.Check(); err != nil || !result {
					fmt.Println(err.Error())
					c.Storage.Del(proxy)
				} else {
					c.Storage.Save(proxy)
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-c.checkTicker.C:
				proxies := c.Storage.All()
				log.Printf("start check from ticker, size %d", len(proxies))
				for _, proxy := range proxies {
					if result, err := proxy.Check(); err != nil || !result {
						fmt.Println(err.Error())
						c.Storage.Del(proxy)
					} else {
						c.Storage.Save(proxy)
					}
				}
			}
		}
	}()
}

func (c *Checker) Stop() {
	if c.checkTicker != nil {
		c.checkTicker.Stop()
	}
	c.done <- true
}
