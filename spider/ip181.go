package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/diinvoke/proxy-pool/model"
	"github.com/diinvoke/proxy-pool/storage"
)

const (
	url = "http://www.ip181.com/"
)

type ip181Result struct {
	Position string `json:"position"`
	Port     string `json:"port"`
	IP       string `json:"ip"`
}

type ip181 struct {
	ErrorCode string         `json:"ERRORCODE"`
	Results   []*ip181Result `json:"RESULT"`
}

type Ip181 struct {
	storage storage.IStorage
	count   int32
}

var _ ISpider = &Ip181{}

func NewIP181(storage storage.IStorage) ISpider {
	return &Ip181{
		storage: storage,
	}
}

func (i *Ip181) Do() error {
	ips, err := i.getIPs()

	i.save(ips)

	return err
}

func (i *Ip181) LoadCount() int32 {
	return atomic.LoadInt32(&i.count)
}

func (i *Ip181) Name() string {
	return "IP181"
}

func (i *Ip181) getIPs() ([]*ip181Result, error) {
	resp, err := http.Get(url)
	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get resp error:%s", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ip181
	if err = json.Unmarshal(b, &result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

func (i *Ip181) save(ips []*ip181Result) {
	if len(ips) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(ips))

	for _, ip := range ips {
		go func(wgg *sync.WaitGroup, ipModel *model.IP) {
			i.storage.Save(ipModel)
			atomic.AddInt32(&i.count, 1)
			wg.Done()
		}(&wg, i.genIPModel(ip))
	}

	wg.Wait()
}

func (i *Ip181) genIPModel(ip *ip181Result) *model.IP {
	return &model.IP{
		Protocol: model.ProtocolHttp,
		Address:  ip.IP,
		Port:     ip.Port,
	}
}
