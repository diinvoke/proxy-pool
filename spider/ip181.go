package spider

import (
	"net/http"
	"strings"
	"time"

	"github.com/Agzdjy/proxy-pool/model"

	"sync"

	"fmt"

	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/Agzdjy/proxy-pool/util"
	"github.com/PuerkitoBio/goquery"
)

type Ip181 struct{}

var _ Spider = &Ip181{}

func (ip181 *Ip181) Do(url string, store storage.Storage) error {
	resp, err := util.HttpGet(url)
	if err != nil {
		return err
	}

	errorChan := make(chan error)
	doneChan := make(chan struct{})

	go func() {
		filterRecord(resp, store, errorChan, doneChan)
	}()

	timeout := time.After(2 * time.Second)

	select {
	case <-doneChan:
		fmt.Println("init ip181 success")
		return nil
	case err := <-errorChan:
		return err
	case <-timeout:
		fmt.Println("init ip181 timeout")
		return nil
	}
	return nil
}

func filterRecord(resp *http.Response, store storage.Storage, errChan chan error, done chan struct{}) {
	var wg sync.WaitGroup

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		errChan <- err
	}

	trs := doc.Find("tbody").Find("tr").Not("tr.active")

	trs.Each(func(index int, tr *goquery.Selection) {
		td := tr.Find("td")
		ipModels := getIpModels(td)

		wg.Add(len(ipModels))
		go checkAndSave(ipModels[0], store, &wg, errChan)
		if len(ipModels) > 1 {
			go checkAndSave(ipModels[1], store, &wg, errChan)
		}
	})

	wg.Wait()
	done <- struct{}{}
}

func getIpModels(td *goquery.Selection) []*model.IP {
	protocol := td.Eq(3).Text()
	protocols := strings.Split(protocol, ",")

	ipModels := []*model.IP{genIpModel(td, protocols[0])}
	if len(protocols) > 1 {
		ipModels = append(ipModels, genIpModel(td, protocols[1]))
	}

	return ipModels
}

func genIpModel(selection *goquery.Selection, protocol string) *model.IP {
	return &model.IP{
		Address:  selection.Eq(0).Text(),
		Port:     selection.Eq(1).Text(),
		Protocol: strings.ToLower(protocol),
	}
}

func checkAndSave(ip *model.IP, store storage.Storage, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	checkUrl := ip.Protocol + "://" + ip.Address + ":" + ip.Port
	if !util.Check(checkUrl) {
		return
	}
	err := store.Save(ip)
	if err != nil {
		errChan <- err
	}
}
