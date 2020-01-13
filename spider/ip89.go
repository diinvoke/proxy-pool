package spider

import (
	"github.com/mingcheng/proxypool/model"
	rpc "github.com/mingcheng/proxypool/protobuf"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type IP89 struct {
}

func (i *IP89) Do() ([]*model.Proxy, error) {
	var ExprIP = regexp.MustCompile(`((25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\:([0-9]+)`)
	pollURL := "http://www.89ip.cn/tqdl.html?api=1&num=100"

	resp, err := http.Get(pollURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyIPs := string(body)

	var results []*model.Proxy
	for _, v := range ExprIP.FindAllString(bodyIPs, 100) {
		address := strings.Split(v, ":")
		if len(address) != 2 {
			continue
		}

		port, err := strconv.ParseUint(address[1], 10, 64)
		if err != nil {
			continue
		}

		results = append(results, &model.Proxy{
			rpc.Proxy{
				Address:  strings.TrimSpace(address[0]),
				Port:     port,
				Protocol: rpc.Protocol_HTTP,
				From:     i.Name(),
			},
		})
	}

	return results, nil
}

func (i *IP89) Name() string {
	return "89ip.cn"
}
