package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mingcheng/proxypool"
	"net/http"
	"time"
)

func main() {

	config := proxypool.Config{
		FetchInterval:   15 * time.Minute,
		CheckInterval:   2 * time.Minute,
		CheckConcurrent: 10,
	}

	go proxypool.Start(config)
	defer proxypool.Stop()

	r := gin.Default()
	r.GET("/all", func(c *gin.Context) {
		proxies := proxypool.All()
		if len(proxies) > 0 {
			c.JSON(http.StatusOK, proxies)
		} else {
			c.String(http.StatusNotFound, "no suitable proxy found")
		}
	})

	r.GET("/random", func(c *gin.Context) {
		if proxies := proxypool.Random(); proxies != nil {
			c.JSON(http.StatusOK, proxies)
		} else {
			c.String(http.StatusNotFound, "no suitable proxy found")
		}
	})

	r.Run()
}
