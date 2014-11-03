package main

import (
	"bloom"
)

func main() {
	bloom := bloom.New(1000)
	bloom.add("foo")
	bloom.IsPresent("foo")
}
