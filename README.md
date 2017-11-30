ip 代理池
=========


### 下载安装

```shell
go get github.com/Agzdjy/proxy-pool
```

### 安装依赖

```shell
glide install
```
[glide](https://github.com/Masterminds/glide)介绍

### 使用

#### 添加存储配置文件config.json

```json
{
  "store": "redis",
  "address": "127.0.0.1:6379",
  "password": "",
  "db": "0"
}
```

#### demo

```go

import (
	"github.com/Agzdjy/proxy-pool"
	"github.com/Agzdjy/proxy-pool/model"
)

func init() {
	proxypool.InitData("./config.json")
}

func GetOne() *model.IP {
	rangeIp := proxypool.Range("https")
	proxyUrl := proxyIp.Protocol + "://" + proxyIp.Address + ":" + proxyIp.Port
	
	if !proxypool.Check(proxyUrl) {
		return GetOne()
	}
	
	return rangeIp
}

func Check(proxy string) bool {
	return proxypool.Check(proxy)
}

```

### 扩展抓取

#### 在``spider``包中编写抓取实现

```go
import (
	"github.com/Agzdjy/proxy-pool/storage"
)

type DemoSpider struct{}

// implements Spider interface
var _ Spider = &DemoSpider{}

func (demo *DemoSpider) Do(url string, store storage.Storage) (count int, err error) {
	// do something
}

```

#### 在``init.go``中初始化抓取

```go
import "github.com/Agzdjy/proxy-pool/spider"

func initSpider(stor storage.Storage) {

	// do something
	&spider.DemoSpider{}.Do(demoUrl, stor)
}
```

### 扩展存储

#### 在``storage``包中编写实现

```go

import (
    "github.com/Agzdjy/proxy-pool/model"
)

type DemoStorage struct{}
// implements Storage interface
var _ Storage = &DemoStorage()

func (d *DemoStorage) Save(ip *model.IP) error {
	// do something
}

func (d *DemoStorage) Del(ip *model.IP) bool {
	// do something
}

func (d *DemoStorage) RangeOne(protocol string) (ip *model.IP, err error) {
	// do something
}

func (d *DemoStorage) Close() error {
	// do something
}


```

#### 在``init.go``中初始化存储

```go

import "github.com/Agzdjy/proxy-pool/storage"

func initStorage(configPath string) storage.Storage {
	config := readJson(configPath)
	store := config["store"]
	var stor storage.Storage

	switch store {
	case "DemoStorage":
		// do something
	}
	return stor
}
```