package main

type Hash struct {
	Rounds     int
	UseIVBlock bool
	InputSize  int
	OutputSize int
	Form       Formula
	w          [64]int //mixing values initialized from input
	chain      []int //
	out        [8]int //output hash value
	in         []int //input byte value (can be variable)
	//a,b,c,d,e,f,g,h mixing values
	A          [70][32]int
	E          [70][32]int
}
