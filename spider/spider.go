package spider

import (
	"github.com/Agzdjy/proxy-pool/storage"
)

type Spider interface {
	Do(url string, store storage.Storage) (count int, err error)
}
