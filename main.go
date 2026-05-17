package main

import (
	"fmt"
)

var (
	debug = false
)

func main() {
	fmt.Println("")
	// test13()
	test17()
	
	initHash("12345678123", 1)
	// listener()
	
}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	h.SetOutput()
	// h.Encode(n)
	writeDIMACS("hash.dimacs", h.F)
}
