package main

import (
	"fmt"
	"crypto/sha256"
)

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

var (
	debug = false
)

func main() {
	x := readDIMACS("out.dimacs")
	// fmt.Println(x)
	writeDIMACS("out2.dimacs", x)
	h := sha256.New()
	v1 := string(h.Sum(readFile("out.dimacs")))
	h.Reset()
	v2 := string(h.Sum(readFile("out2.dimacs")))
	if v1 == v2 {
		fmt.Println("Matches", v1)
	}
}

type Hash struct {
	Rounds     int
	UseIVBlock bool
	InputSize  int
	OutputSize int
	Form       Formula
	w          []byte
	chain      []byte
	out        []byte
	in         []byte
}
