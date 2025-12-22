package main

import (
	"fmt"
)

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

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
	// fmt.Println(h.F)
	h.Encode()
	writeDIMACS("hash.dimacs", h.F)
}
