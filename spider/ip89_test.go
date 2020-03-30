package spider

import (
	"testing"
)

func TestIP89_Do(t *testing.T) {
	i := &IP89{}
	if got, err := i.Do(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("got %d proxies from %s", len(got), i.Name())
	}
}
