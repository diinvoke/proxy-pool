package spider

import (
	"net/http"
	"strings"

	"github.com/Agzdjy/go-proxy/model"

	"github.com/Agzdjy/go-proxy/storage"
	"github.com/PuerkitoBio/goquery"
)

type ip181 struct{}

var _ Spider = &ip181{}

func (ip181 *ip181) Do(url string, store storage.Storage) (count int, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

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

	trs := doc.Find("tbody").Find("tr.warning")

	trs.Each(func(index int, tr *goquery.Selection) {
		td := tr.Find("td")
		protocolStr := td.Eq(3).Text()
		protocols := strings.Split(protocolStr, ",")

		store.Save(&model.IP{
			Address:  td.Eq(0).Text(),
			Port:     td.Eq(1).Text(),
			Protocol: protocols[0],
		})

		if len(protocols) > 1 {
			store.Save(&model.IP{
				Address:  td.Eq(0).Text(),
				Port:     td.Eq(1).Text(),
				Protocol: protocols[1],
			})
		}

		count++
	})

	return count, err
}
