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

/* Encode sha256 as CNF */
func (h *Hash) Encode() {
	//for len(h.P) >= 64 {
	h.InitVars()
	for j := 0; j < 1; j++ {
		//todo: chunked
		h.EncodeConst()
		/* need to verify */
		if h.UseIVBlock {
			h.EncodeIV()
			h.UseIVBlock = false
		}
		/* need to verify */
		for i := 0; i < 16; i++ {
			h.Mixer1(i)
		}
		for i := 16; i < 64; i++ {
			h.Mixer2(i)
		}
		// h.EncodeConst()
		/* need to verify */
		// if h.UseIVBlock {
		// 	h.EncodeIV()
		// 	h.UseIVBlock = false
		// }
		h.AssignChain()
		//Main loop
		for i := 0; i < 64; i++ {
			h.MainLoop(i)
		}
		//final addition
		h.Add()
		//p = p[chunk:]
		h.ChunkRounds += 1
	}
	h.UseIVBlock = true
	//assign to output
	h.SetOutput()
	h.ChunkRounds = 0
	h.EncodeRounds += 1
}

/* Set the output bits and hash output string */
func (h *Hash) SetOutput() {
	for i := 0; i < 8; i++ {
		h.F.Assign(arrConv(h.out[i]), arrConv(h.tmp[i]), 32)
	}
	h.Output = ""
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
	sig0 := h.F.NewVars(32, fmt.Sprintf("sig0_%d", i))
	sig1 := h.F.NewVars(32, fmt.Sprintf("sig1_%d", i))
	h.F.Sigma0(sig0, arrConv(h.A[i+3]))
	h.F.Sigma1(sig1, arrConv(h.E[i+3]))
	f1 := h.F.NewVars(32, fmt.Sprintf("f1_%d", i))
	f2 := h.F.NewVars(32, fmt.Sprintf("f2_%d", i))
	h.F.If(f1, arrConv(h.E[i+3]), arrConv(h.E[i+2]), arrConv(h.E[i+1]), 32)
	h.F.Maj3(f2, arrConv(h.A[i+3]), arrConv(h.A[i+2]), arrConv(h.A[i+1]), 32)
	t := h.F.NewVars(32, fmt.Sprintf("t_%d", i))
	h.F.Add5(t, arrConv(h.E[i]), sig1, f1, arrConv(h.K[i]), arrConv(h.W[i]), 32)
	h.F.Add2(arrConv(h.E[i+4]), arrConv(h.A[i]), t, 32)
	h.F.Add3(arrConv(h.A[i+4]), t, sig0, f2, 32)
}

/* Mix p values into w, padded as big endian */
func (h *Hash) Mixer1(i int) {
	j := i * 4
	_ = j
	//w[i] = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
	// x := uint32(h.P[j])<<24 | uint32(h.P[j+1])<<16 | uint32(h.P[j+2])<<8 | uint32(h.P[j+3])
	// h.W[i] = [32]int(intToBitVec(x))
	// h.W[i] = [32]int(h.F.Assign(arrConv(h.W[i]), intToBitVec(uint32(bitVecToInt(arrConv(h.P[j])))<<24 | uint32(bitVecToInt(arrConv(h.P[j+1])))<<16 | uint32(bitVecToInt(arrConv(h.P[j+2])))<<8 | uint32(bitVecToInt(arrConv(h.P[j+3]))), 32)))

	// h.W[i] = h.P[i]
}

/* Mix w values */
func (h *Hash) Mixer2(i int) {
	s0 := h.F.NewVars(32, "s0")
	s1 := h.F.NewVars(32, "s1")
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
	h.F.Assign(arrConv(h.A[3]), arrConv(h.chain[0]), 32)
	h.F.Assign(arrConv(h.A[2]), arrConv(h.chain[1]), 32)
	h.F.Assign(arrConv(h.A[1]), arrConv(h.chain[2]), 32)
	h.F.Assign(arrConv(h.A[0]), arrConv(h.chain[3]), 32)
	h.F.Assign(arrConv(h.E[3]), arrConv(h.chain[4]), 32)
	h.F.Assign(arrConv(h.E[2]), arrConv(h.chain[5]), 32)
	h.F.Assign(arrConv(h.E[1]), arrConv(h.chain[6]), 32)
	h.F.Assign(arrConv(h.E[0]), arrConv(h.chain[7]), 32)
}

/* Encode IV constants */
func (h *Hash) EncodeIV() {
	for i := 0; i < 8; i++ {
		h.F.FixedValue(arrConv(h.chain[i]), int(_IV[i]), 32)
	}
}

/* Encode k values */
func (h *Hash) EncodeConst() {
	for i := 0; i < 64; i++ {
		h.F.NewVars(len(h.K[i]), fmt.Sprintf("k_%d", i))
		h.F.FixedValue(arrConv(h.K[i]), int(_K[i]), 32)
	}
}

/* Init (chunk dependent) named cnf variables */
func (h *Hash) InitVars() {
	for i := 0; i < 8; i++ {
		h.chain[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("chain_%d", i)))
		h.tmp[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("tmp_%d", i)))
	}
	for i := 0; i < h.Rounds; i++ {
		h.W[i] = [32]int(h.F.NewVars(32, fmt.Sprintf("w_%d", i)))
		h.A[i+4] = [32]int(h.F.NewVars(32, fmt.Sprintf("A_%d", i+4)))
		h.E[i+4] = [32]int(h.F.NewVars(32, fmt.Sprintf("E_%d", i+4)))
	}
}

/* Chained hashes */
func (h *Hash) MultiHash() {
	for i := 0; i < len(h.Table); i++ {
		h.Encode()
		h.Input = h.Output
		//output vars moved to input
		h.InitBitVec()
	}
}

/* Create a new sha256 instance */
func NewSHA256(input string, n int) Hash {
	h := Hash{}
	h.Rounds = 64
	h.F = NewFormula()
	h.EncodeRounds = 0
	h.ChunkRounds = 0
	h.Input = input
	h.Table = make([]VecTable, n)
	h.UseIVBlock = true
	h.InitBitVec()
	// h.OutputSize = 8
	return h
}

/* Init the hashes bit vectors */
func (h *Hash) InitBitVec() {
	//todo: needs padding
	h.P = make([][8]int, len(h.Input))
	for i := 0; i < len(h.Input); i++ {
		h.P[i] = [8]int(byteToBitBec(h.Input[i]))
	}
	h.InitVars()
	// h.InputSize = len(h.Input)*8
}

/*
func (d *digest) Write(p []byte) (nn int, err error) {
	boring.Unreachable()
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		d.nx += n
		if d.nx == chunk {
			blockGeneric(d, d.x[:])
			d.nx = 0
		}
		p = p[n:]
	}
	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		blockGeneric(d, p[:n])
		p = p[n:]
	}
	if len(p) > 0 {
		d.nx = copy(d.x[:], p)
	}
	return
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	h     [8]uint32
	x     [chunk]byte
	nx    int
	len   uint64
}
// The size of a SHA256 checksum in bytes.
const Size = 32
// The blocksize of SHA256 and SHA224 in bytes.
const BlockSize = 64
// Sum256 returns the SHA256 checksum of the data.
func Sum256(data []byte) [Size]byte {
	if boring.Enabled {
		return boring.SHA256(data)
	}
	var d digest
	d.Reset()
	d.Write(data)
	return d.checkSum()
}
*/

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
    copy chunk into first 16 words w[0..15] of the message schedule array
*/

type Hash struct {
	F            Formula
	Table        []VecTable //bitvec table
	Multi        int        //number of chained hashes
	Rounds       int        //64
	EncodeRounds int        //
	ChunkRounds  int        //
	UseIVBlock   bool       //
	InputSize    int        //in bits (always a multiple of the number of bits in a character)
	OutputSize   int        //
	Input        string     //raw input
	Output       string     //output string
	chain        [8][32]int //
	tmp          [8][32]int //
	out          [8][32]int //output hash value
	P            [][8]int   //input string as bitvec, must be padded
	//a,b,c,d,e,f,g,h mixing values
	A [70][32]int
	E [70][32]int
	K [64][32]int //k constants
	W [64][32]int //mixing values initialized from input
}

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
