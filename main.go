package main

import (
	"fmt"
)

var (
	debug  = true
	mixer1 = false
)

func main() {
	// test13()
	// test7()
	// test17()
	test19()
	// return
	var x string
	if mixer1 {
		fmt.Scanf("%s", x)
		initHash(x, 1)
	} else {
		initHash("123", 1)
	}
	// listener()
	// r := sha256sum("123")
	// e := stringToBitVec(sha256sum(("123")))
	// fmt.Println(e)
	// r = fmt.Sprintf("%x", r)
	// fmt.Println(r)
	

}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	// h.Encode(n)
	h.MultiHash()
	writeDIMACS("hash.dimacs", h.F)
}







