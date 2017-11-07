package spider

import (
	"net/http"
	"strings"

	"github.com/Agzdjy/proxy-pool/model"

	"github.com/Agzdjy/proxy-pool/storage"
	"github.com/PuerkitoBio/goquery"
)

type Ip181 struct{}

var _ Spider = &Ip181{}

func (ip181 *Ip181) Do(url string, store storage.Storage) (count int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return filterRecord(resp, store)
}

func filterRecord(response *http.Response, store storage.Storage) (count int, err error) {
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return 0, err
	}

	trs := doc.Find("tbody").Find("tr").Not("tr.active")

	done := make(chan string, 10)
	stop := make(chan struct{})

	trs.Each(func(index int, tr *goquery.Selection) {
		td := tr.Find("td")
		protocolStr := td.Eq(3).Text()
		protocols := strings.Split(protocolStr, ",")

		httpIp := &model.IP{
			Address:  td.Eq(0).Text(),
			Port:     td.Eq(1).Text(),
			Protocol: strings.ToLower(protocols[0]),
		}
		go saveToStorage(httpIp, store, stop, done)

		if len(protocols) > 1 {
			httpsIp := &model.IP{
				Address:  td.Eq(0).Text(),
				Port:     td.Eq(1).Text(),
				Protocol: strings.ToLower(protocols[1]),
			}
			go saveToStorage(httpsIp, store, stop, done)
			count++
		}

		count++
	})

	for i := 0; i < count; i += 1 {
		<-done
	}
	close(stop)
	return count, err
}

func saveToStorage(ip *model.IP, store storage.Storage, stop <-chan struct{}, done chan string) error {
	err := store.Save(ip)

	select {
	case <-stop:
		_, isClose := <-done
		if !isClose {
			close(done)
		}
		return nil
	case done <- "save ok":
	}
	return err
}
