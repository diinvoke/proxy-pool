package spider

import (
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/mingcheng/proxypool/model"
	rpc "github.com/mingcheng/proxypool/protobuf"
)

type Feiyi struct {
	// URL   string
	// Query string
}

func (i *Feiyi) Do() ([]*model.Proxy, error) {
	const (
		_URL   = "http://www.feiyiproxy.com/?page_id=1457"
		_Query = "//div[contains(@class, 'et_pb_code_1')]//tr"
	)

	doc, err := htmlquery.LoadURL(_URL)
	if err != nil {
		return nil, err
	}

	trNodes, err := htmlquery.QueryAll(doc, _Query)
	if err != nil {
		return nil, err
	}

	var results []*model.Proxy
	for _, trNode := range trNodes {
		tdNode := htmlquery.Find(trNode, "//td")
		if tdNode == nil || len(tdNode) <= 0 {
			continue
		}

		port, err := strconv.ParseUint(htmlquery.InnerText(tdNode[1]), 10, 64)
		if err != nil {
			continue
		}

		results = append(results, &model.Proxy{
			rpc.Proxy{
				Address:  strings.TrimSpace(htmlquery.InnerText(tdNode[0])),
				Port:     port,
				Protocol: model.ProtocolFromString(htmlquery.InnerText(tdNode[3])),
				From:     i.Name(),
			},
		})
	}

	return results, nil
}

func (i *Feiyi) Name() string {
	return "feiyiproxy.com"
}
