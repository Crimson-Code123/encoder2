package main

import (
	"fmt"
)

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

var (
	debug = false
)

func main() {
	fmt.Println("")
	test5()
	listener()
}

func initHash(input string) {
	h := NewSHA256(input)
	h.Encode()
	writeDIMACS("hash.dimacs", h.F)
}



