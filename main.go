package main

import (
	"fmt"
)

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

var (
	
)


func main() {
	x := readDIMACS("out.dimacs")
	fmt.Println(x)

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


