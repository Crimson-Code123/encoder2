package main

/* sha256 cnf-encoder */

/* Encode sha256 as CNF */
func (h *Hash) Encode() {
	//todo: chunked
	/* need to verify */
	for i := 0; i < 16; i++ {
		h.Mixer1(i)
	}
	for i := 16; i < 64; i++ {
		h.Mixer2(i)
	}
	h.EncodeConst()
	/* need to verify */
	if h.UseIVBlock {
		h.EncodeIV()
	}
	h.AssignChain()
	//Main loop
	for i := 0; i < 64; i++ {
		h.MainLoop(i)
	}
	//final addition ()
	h.Add()
}

/* Final addition in sha256 loop (for readability) */
func (h *Hash) Add() {
	/*
		cnf.add2(out[0], chain[0], A[rounds+3]);
		cnf.add2(out[1], chain[1], A[rounds+2]);
		cnf.add2(out[2], chain[2], A[rounds+1]);
		cnf.add2(out[3], chain[3], A[rounds]);
		cnf.add2(out[4], chain[4], E[rounds+3]);
		cnf.add2(out[5], chain[5], E[rounds+2]);
		cnf.add2(out[6], chain[6], E[rounds+1]);
		cnf.add2(out[7], chain[7], E[rounds]);
	*/
	// h.F.Add2(h.out[0], h.chain[0], h.A[64+3], 32)
	// h.F.Add2()
	// h.F.Add2()
	// h.F.Add2()
	// h.F.Add2()
	// h.F.Add2()
	// h.F.Add2()
	// h.F.Add2()
}

func (h *Hash) MainLoop(i int) {
	sig0 := h.F.NewVars(32, "sig0")
	sig1 := h.F.NewVars(32, "sig1")
	_ = sig0
	_ = sig1
	// h.F.Sigma0(sig0)
	// h.F.Sigma1(sig1)
	// f1 := h.F.NewVars(32, "f1")
	// f2 := h.F.NewVars(32, "f2")
	// h.F.If(f1, h.E[i+3], h.E[i+2], h.E[i+1], 32)
	// h.F.Maj3(f2, h.A[i+3], h.A[i+2], h.A[i+1], 32)
	// t := h.F.NewVars(32, "t")
	// h.F.Add5(t, h.E[:], sig1, f1, _K[i], h.w[i], 32)
	// h.F.Add2(h.E[:], h.A[:], t)
	// h.F.Add3(h.A[:], T, sig0, f2)
}

func (h *Hash) Mixer1(i int) {

}

/* Mix w values */
func (h *Hash) Mixer2(i int) {
	s0 := h.F.NewVars(32, "s0")
	s1 := h.F.NewVars(32, "s1")
	_ = s0
	_ = s1
	// r1 := h.F.RotateR(h.w[i-15], 7)
}

/* Assign chaining values to a,b,c,d,e,f,g,h */
func (h *Hash) AssignChain() {
	/*
		cnf.assign(A[3], chain[0]);
		cnf.assign(A[2], chain[1]);
		cnf.assign(A[1], chain[2]);
		cnf.assign(A[0], chain[3]);
		cnf.assign(E[3], chain[4]);
		cnf.assign(E[2], chain[5]);
		cnf.assign(E[1], chain[6]);
		cnf.assign(E[0], chain[7]);
	*/

	// h.F.Assign(h.A[3], h.chain[0], 32)
	// h.F.Assign()
	// h.F.Assign()
	// h.F.Assign()
	// h.F.Assign()
	// h.F.Assign()
	// h.F.Assign()
	// h.F.Assign()
}

/*  */
func (h *Hash) EncodeIV() {
	/*
		cnf.fixedValue(chain[0], 0x6a09e667);
		cnf.fixedValue(chain[1], 0xbb67ae85);
		cnf.fixedValue(chain[2], 0x3c6ef372);
		cnf.fixedValue(chain[3], 0xa54ff53a);
		cnf.fixedValue(chain[4], 0x510e527f);
		cnf.fixedValue(chain[5], 0x9b05688c);
		cnf.fixedValue(chain[6], 0x1f83d9ab);
		cnf.fixedValue(chain[7], 0x5be0cd19);
	*/
	for i := 0; i < 8; i++ {
		// h.F.FixedValue(h.chain[i], int(_IV[i]), 32)
	}
}

/*  */
func (h *Hash) EncodeConst() {
	var k = [64][32]int{}
	for i := 0; i < 64; i++ {
		h.F.NewVars(len(k[i]), "k")
		h.F.FixedValue(arrayConv(k[i]), int(_K[i]), 32)
	}
}

func NewSHA256(input string) Hash {
	h := Hash{}
	h.Rounds = 64
	h.in = []byte(input)
	h.F = NewFormula()
	h.UseIVBlock = true
	h.in = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		h.in[i] = byte(input[i])
	}
	h.InputSize = len(input)
	h.OutputSize = 8
	return h
}

var _K = []uint32{
	0x428a2f98,
	0x71374491,
	0xb5c0fbcf,
	0xe9b5dba5,
	0x3956c25b,
	0x59f111f1,
	0x923f82a4,
	0xab1c5ed5,
	0xd807aa98,
	0x12835b01,
	0x243185be,
	0x550c7dc3,
	0x72be5d74,
	0x80deb1fe,
	0x9bdc06a7,
	0xc19bf174,
	0xe49b69c1,
	0xefbe4786,
	0x0fc19dc6,
	0x240ca1cc,
	0x2de92c6f,
	0x4a7484aa,
	0x5cb0a9dc,
	0x76f988da,
	0x983e5152,
	0xa831c66d,
	0xb00327c8,
	0xbf597fc7,
	0xc6e00bf3,
	0xd5a79147,
	0x06ca6351,
	0x14292967,
	0x27b70a85,
	0x2e1b2138,
	0x4d2c6dfc,
	0x53380d13,
	0x650a7354,
	0x766a0abb,
	0x81c2c92e,
	0x92722c85,
	0xa2bfe8a1,
	0xa81a664b,
	0xc24b8b70,
	0xc76c51a3,
	0xd192e819,
	0xd6990624,
	0xf40e3585,
	0x106aa070,
	0x19a4c116,
	0x1e376c08,
	0x2748774c,
	0x34b0bcb5,
	0x391c0cb3,
	0x4ed8aa4a,
	0x5b9cca4f,
	0x682e6ff3,
	0x748f82ee,
	0x78a5636f,
	0x84c87814,
	0x8cc70208,
	0x90befffa,
	0xa4506ceb,
	0xbef9a3f7,
	0xc67178f2,
}

var _IV = []uint32{
	0x6a09e667,
	0xbb67ae85,
	0x3c6ef372,
	0xa54ff53a,
	0x510e527f,
	0x9b05688c,
	0x1f83d9ab,
	0x5be0cd19,
}

type Hash struct {
	F          Formula
	Rounds     int
	UseIVBlock bool
	InputSize  int
	OutputSize int
	w          [64]int //mixing values initialized from input
	chain      []int   //
	out        [8]int  //output hash value
	in         []byte  //input byte value (can be variable)
	//a,b,c,d,e,f,g,h mixing values
	A [70][32]int
	E [70][32]int
}
