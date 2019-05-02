ip 代理池
=========


### 下载安装

```shell
go get github.com/diinvoke/proxy-pool
```

### 安装依赖

```shell
dep ensure
```
[dep](https://golang.github.io/dep/docs/introduction.html)介绍

### 使用

#### demo

```go

import (
	"github.com/diinvoke/proxy-pool"
	"github.com/diinvoke/proxy-pool/model"
)

func Random() *model.IP {
	rangeIp := proxypool.Range("https")
	
	// if need check
	// proxyUrl := proxyIp.Protocol + "://" + proxyIp.Address + ":" + proxyIp.Port
	// if !proxypool.Check(proxyUrl) {
		// return Random()
	// }
	
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
	"github.com/diinvoke/proxy-pool/storage"
)

type DemoSpider struct{}

func NewDemoSpider() ISpider {
	return &DemoSpider{}
}

// implements Spider interface
var _ ISpider = &DemoSpider{}

func (demo *DemoSpider) Do() error {
	// do something
}

```

#### 在``init.go``中添加或者替换实现

```go
import "github.com/diinvoke/proxy-pool/spider"

spiders := []spider.ISpider{
	spider.NewIP181(store),
	spider.NewDemoSpider()
}
```

### 扩展存储

#### 在``storage``包中编写实现

```go

import (
    "github.com/diinvoke/proxy-pool/model"
)

type DemoStorage struct{}
// implements Storage interface
var _ IStorage = &DemoStorage()

func NewDemoStorage() IStorage {
	return &DemoStorage{}
}

func (d *DemoStorage) Save(ip *model.IP) error {
	// do something
}

func (d *DemoStorage) Del(ip *model.IP) bool {
	// do something
}

func (d *DemoStorage) Random(protocol model.Protocol) (ip *model.IP, err error) {
	// do something
}

func (d *DemoStorage) Close() error {
	// do something
}


```

#### 替换 ``storage/storage.go GetStorage 返回值``

```go
func GetStorage() IStorage {
	return NewDemoStorage()
}
```
