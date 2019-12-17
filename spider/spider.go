package spider

import (
	"github.com/mingcheng/proxypool/model"
)

type Spider interface {
	Do() ([]*model.Proxy, error)
	Name() string
}
