package spider

import (
	"testing"
)

func TestFeiyi_Do(t *testing.T) {
	i := &Feiyi{}
	if got, err := i.Do(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("got %d proxies from %s", len(got), i.Name())
	}
}
