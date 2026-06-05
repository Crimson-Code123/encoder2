package main

import (
	"fmt"
)

var (
	debug = false
	mixer1 = false
	
)

/*  */

func main() {
	fmt.Printf("")
	// test13()
	// test7()
	
	initHash("123", 1)
	// listener()
	
}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	// h.SetOutput()
	// h.Encode(n)
	h.MultiHash()
	writeDIMACS("hash.dimacs", h.F)
}




