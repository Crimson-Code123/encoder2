package main

import (
	"fmt"
)

var (
	debug = false
)

func main() {
	fmt.Sprintf("")
	// initHash("12345678123", 1)
	// listener()
	
}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	h.Encode()
	writeDIMACS("hash.dimacs", h.F)
}
