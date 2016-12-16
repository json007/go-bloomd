package bloomfilter

import (
	"test"
)

func Testadd(t *testing.T) {
	bf := new(100)
	bf.Add([]byte("test"))
	if !bf.IsContain([]byte("test")) {
		t.Errorf("add and check error")
	}
}
