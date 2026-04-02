package main

import "fmt"

/* sha256 cnf-encoder */

/*
  Note 1: All variables are 32 bit unsigned integers (and addition is calculated modulo 2^32)
  Note 2: For each round, there is one round constant k[i] and one entry in the message schedule array w[i], 0 ≤ i ≤ 63
  Note 3: The compression function uses 8 working variables, a through h
  Note 4: Big-endian convention is used when expressing the constants in this pseudocode,
  	and when parsing message block data from bytes to words, for example,
  	the first word of the input message "abc" after padding is 0x61626380
*/

/* todo
functions -
	Mixer1 -
	initvars -
	setinput -
	initbitvec -
	Encode - chunking
	Setoutput -
*/

/*

 */

/* test values
empty input
2080933965 1779033703 3144134277 1013904242 415752866 1359893119 2600822924 528734635 - abc values, 1 round
3820012610 2566659092 2600203464 2574235940 665731556 1687917388 2761267483 2018687061 - digest output
[227 176 196 66 152 252 28 20 154 251 244 200 153 111 185 36 39 174 65 228 100 155 147 76 164 149 153 27 120 82 184 85] - checksum
*/

/* Encode sha256 as CNF */
func (h *Hash) Encode(j int) {
	h.UseIVBlock = true
	// for len(h.P) >= 64 {

	h.InitVars(j)
	//h.SetInput()
	h.EncodeConst(j)
	if h.UseIVBlock {
		h.EncodeIV()
		h.UseIVBlock = false
	}
	// for len(h.P) >= h.ChunkSize {
	for x := 0; x < 1; x++ {
		//todo: chunked
		for i := 0; i < 16; i++ {
			/*
				this portion should not require encoding
				it is required for fixing the input bits, however
			*/
			h.Mixer1(i)
		}
		for i := 16; i < 64; i++ {
			h.Mixer2(i)
		}
		h.AssignChain()
		for i := 0; i < 64; i++ {
			h.MainLoop(i)
		}
		//final addition
		h.Add()
		// h.P = h.P[h.ChunkSize:]
	}
	//assign to output
	h.SetOutput()
	// h.SetInput(h.out)
}

/* Set the output bits and hash output string */
func (h *Hash) SetOutput() {
	for i := 0; i < 8; i++ {
		h.out[i] = [32]int(h.F.Assign(arrConv(h.out[i]), arrConv(h.tmp[i]), 32))
		//todo: fix output bits (get from btc difficulty)
		h.F.FixedValue(arrConv(h.out[i]), uint32(0), 32)
	}
	// h.Output = bitVecToString(h.out)
}

/* Final addition in sha256 loop */
func (h *Hash) Add() {
	h.F.Add2(arrConv(h.tmp[0]), arrConv(h.chain[0]), arrConv(h.A[64+3]), 32)
	h.F.Add2(arrConv(h.tmp[1]), arrConv(h.chain[1]), arrConv(h.A[64+2]), 32)
	h.F.Add2(arrConv(h.tmp[2]), arrConv(h.chain[2]), arrConv(h.A[64+1]), 32)
	h.F.Add2(arrConv(h.tmp[3]), arrConv(h.chain[3]), arrConv(h.A[64]), 32)
	h.F.Add2(arrConv(h.tmp[4]), arrConv(h.chain[4]), arrConv(h.E[64+3]), 32)
	h.F.Add2(arrConv(h.tmp[5]), arrConv(h.chain[5]), arrConv(h.E[64+2]), 32)
	h.F.Add2(arrConv(h.tmp[6]), arrConv(h.chain[6]), arrConv(h.E[64+1]), 32)
	h.F.Add2(arrConv(h.tmp[7]), arrConv(h.chain[7]), arrConv(h.E[64]), 32)
}

/* main loop of sha256 */
func (h *Hash) MainLoop(i int) {
	sig0 := h.F.NewVars(32, "" /* fmt.Sprintf("SIG0_%d", i) */)
	sig1 := h.F.NewVars(32, "" /* fmt.Sprintf("SIG1_%d", i) */)
	h.F.Sigma0(sig0, arrConv(h.A[i+3]))
	h.F.Sigma1(sig1, arrConv(h.E[i+3]))
	f1 := h.F.NewVars(32, "" /* fmt.Sprintf("F1_%d", i) */)
	f2 := h.F.NewVars(32, "" /* fmt.Sprintf("F2_%d", i) */)
	h.F.If(f1, arrConv(h.E[i+3]), arrConv(h.E[i+2]), arrConv(h.E[i+1]), 32)
	/* c++ - 32 bit diff
	-19409 2049 2081 0
	-19409 2049 2113 0
	-19409 2081 2113 0
	19409 -2081 -2113 0
	19409 -2049 -2113 0
	19409 -2049 -2081 0
	go - 64 bit diff
	-26065 6657 6721 0
	-26065 6657 6785 0
	-26065 6721 6785 0
	26065 -6721 -6785 0
	26065 -6657 -6785 0
	26065 -6657 -6721 0
	*/
	h.F.Maj3(f2, arrConv(h.A[i+3]), arrConv(h.A[i+2]), arrConv(h.A[i+1]), 32)
	t := h.F.NewVars(32, "" /* fmt.Sprintf("T_%d", i) */)
	h.F.Add5(t, arrConv(h.E[i]), sig1, f1, arrConv(h.K[i]), arrConv(h.W[i]), 32)
	h.F.Add2(arrConv(h.E[i+4]), arrConv(h.A[i]), t, 32)
	h.F.Add3(arrConv(h.A[i+4]), t, sig0, f2, 32)
}

/* Mix p values into w, padded as big endian */
func (h *Hash) Mixer1(i int) {

	j := i * 4
	x := bitVecToInt(arrConv(h.P[j]))<<24 | bitVecToInt(arrConv(h.P[j+1]))<<16 | bitVecToInt(arrConv(h.P[j+2]))<<8 | bitVecToInt(arrConv(h.P[j+3]))
	h.W[i] = [32]int(intToBitVec(x))
	//todo: fix the input bits
	//h.F.FixedValue()
	
}

/* Fix the input bits and pad the input string into a bitvec */
func (h *Hash) SetInput(s string) {
	h.FixInput()
	//todo: set h input from s
	q := len(h.Input) % 4
	// fmt.Println(q)
	if q != 0 {
		h.Input += padString(4 - q - 1)
	}
	// fmt.Println(h.Input)
	q = len(h.Input) / 4
	h.P = make([][32]int, len(h.Input))
	// quit()
	for i := 0; i < len(h.Input); i += 4 {
		// h.P[i] = [32]int(byteToBitBec(h.Input[i]))
		// h.P[i] = shaPaddedMsg(h.Input[i : i+4])
	}
}

/* Mix w values */
func (h *Hash) Mixer2(i int) {
	s0 := h.F.NewVars(32, "" /* fmt.Sprintf("S0_%d", i) */)
	s1 := h.F.NewVars(32, "" /* fmt.Sprintf("S1_%d", i) */)
	r1 := make([]int, 32)
	r2 := make([]int, 32)

	r1 = h.F.RotateR(r1, arrConv(h.W[i-15]), 7, 32)
	r2 = h.F.RotateR(r2, arrConv(h.W[i-15]), 18, 32)
	h.F.Xor2(s0[29:], r1[29:], r2[29:], 3)
	h.F.Xor3(s0, r1, r2, h.W[i-15][3:], 29) //shr(w[i-15]+3)

	r1 = h.F.RotateR(r1, arrConv(h.W[i-2]), 17, 32)
	r2 = h.F.RotateR(r2, arrConv(h.W[i-2]), 19, 32)

	h.F.Xor2(s1[22:], r1[22:], r2[22:], 10)
	h.F.Xor3(s1, r1, r2, h.W[i-2][10:], 22) //shr(w[i-2]+10)

	h.F.Add4(arrConv(h.W[i]), arrConv(h.W[i-16]), s0, arrConv(h.W[i-7]), s1, 32)
}

/* Assign chaining values to a,b,c,d,e,f,g,h */
func (h *Hash) AssignChain() {
	h.A[3] = [32]int(h.F.Assign(arrConv(h.A[3]), arrConv(h.chain[0]), 32))
	h.A[2] = [32]int(h.F.Assign(arrConv(h.A[2]), arrConv(h.chain[1]), 32))
	h.A[1] = [32]int(h.F.Assign(arrConv(h.A[1]), arrConv(h.chain[2]), 32))
	h.A[0] = [32]int(h.F.Assign(arrConv(h.A[0]), arrConv(h.chain[3]), 32))
	h.E[3] = [32]int(h.F.Assign(arrConv(h.E[3]), arrConv(h.chain[4]), 32))
	h.E[2] = [32]int(h.F.Assign(arrConv(h.E[2]), arrConv(h.chain[5]), 32))
	h.E[1] = [32]int(h.F.Assign(arrConv(h.E[1]), arrConv(h.chain[6]), 32))
	h.E[0] = [32]int(h.F.Assign(arrConv(h.E[0]), arrConv(h.chain[7]), 32))
}

/* Encode IV constants */
func (h *Hash) EncodeIV() {
	for i := 0; i < 8; i++ {
		h.F.FixedValue(arrConv(h.chain[i]), _IV[i], 32)
	}
}

/* Encode k values */
func (h *Hash) EncodeConst(j int) {
	for i := 0; i < 64; i++ {
		h.K[i] = [32]int(h.F.NewVars(len(h.K[i]), fmt.Sprintf("k%d%d", i, j)))
		h.F.FixedValue(arrConv(h.K[i]), _K[i], 32)
	}
}

/* Init (chunk dependent) named cnf variables */
func (h *Hash) InitVars(d int) {
	for i := 0; i < h.Rounds; i++ {
		h.W[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("w%d%d", i, d)))
	}
	for i := 0; i < 8; i++ {
		h.chain[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("chain%d%d", i, d)))
	}
	for i := 0; i < 8; i++ {
		h.tmp[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("tmp%d%d", i, d)))
	}
	for i := 0; i < h.Rounds; i++ {
		h.A[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("a%d%d", i, d)))
		h.E[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("e%d%d", i, d)))
	}
}

/* Chained hashes */
func (h *Hash) MultiHash() {
	h.SetInput(h.Input)
	for i := 0; i < len(h.Table); i++ {
		
		h.Encode(i)
		h.SetInput(bitVecToString((h.out)))
		h.Input = h.Output
		//output vars moved to input
		// h.InitBitVec()
	}
}

/* Create a new sha256 instance */
func NewSHA256(input string, n int) Hash {
	h := Hash{}
	h.Rounds = 64
	h.ChunkSize = 16
	h.FixedBits = 0 //
	h.F = NewFormula()
	h.Input = input
	h.Table = make([]VecTable, n)
	h.UseIVBlock = true
	// h.InitBitVec()
	return h
}

/* Init the hashes bit vectors */
func (h *Hash) InitBitVec() {
	h.SetInput("")
	// h.InitVars()
	h.InputSize = len(h.Input) * 8
}

/* Fix input bits (mixing values) */
func (h *Hash) FixInput() {
	for i := 0; i < h.FixedBits; i++ {
		r := i / 32 //int size
		c := i % 32
		x := bitVecToInt(arrConv(h.W[r]))
		h.F.FixedValue(intToBitVec(uint32(h.W[r][c])), x>>c&1, 1)
	}
}

/* first 32 bits of the fractional parts of the cube roots of the first 64 primes 2..311) */
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

/* first 32 bits of the fractional parts of the square roots of the first 8 primes 2..19):*/
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

/*
    create a 64-entry message schedule array w[0..63] of 32-bit words
	(The initial values in w[0..63] don't matter, so many implementations zero them here)
    copy chunk into first 16 words w[0..15] of the message schedule array (512 bits)
*/

type Hash struct {
	F          Formula
	Table      []VecTable //bitvec table
	Multi      int        //number of chained hashes
	Rounds     int        //64
	UseIVBlock bool       //
	InputSize  int        //in bits (always a multiple of the number of bits in a character)
	FixedBits  int        //number of fixed input bits
	OutputSize int        //
	Input      string     //raw input
	Output     string     //output string
	ChunkSize  int        //chunk size in bytes
	chain      [8][32]int //
	tmp        [8][32]int //
	out        [8][32]int //output hash value
	P          [][32]int  //input string as bitvec, must be padded
	//a,b,c,d,e,f,g,h mixing values
	A [70][32]int
	E [70][32]int
	K [64][32]int //k constants
	W [64][32]int //mixing values initialized from input
}

/* Bit vector table for storing multi hash values */
type VecTable struct {
	chain [8][32]int //
	tmp   [8][32]int //
	out   [8][32]int //output hash value
	P     [][8]int   //input string as bitvec, must be padded
	//a,b,c,d,e,f,g,h mixing values
	A [70][32]int
	E [70][32]int
	K [64][32]int //k constants
	W [64][32]int //mixing values initialized from input
}
