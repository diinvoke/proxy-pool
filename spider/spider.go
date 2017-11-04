package spider

import (
	"github.com/Agzdjy/go-proxy/storage"
)

type Spider interface {
	Do(url string, store storage.Storage) (count int, err error)
}
