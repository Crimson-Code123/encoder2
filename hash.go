package main

type Hash struct {
	Rounds     int
	UseIVBlock bool
	InputSize  int
	OutputSize int
	Form       Formula
	w          [64]byte //mixing values initialized from input
	chain      []byte //
	out        [8]byte //output hash value
	in         []byte //input byte value (can be variable)
	//a,b,c,d,e,f,g,h mixing values
	A          [70]byte
	E          [70]byte
}
