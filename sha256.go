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
	for j := 0; j < 1; j++ {
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
		//final addition
		h.Add()
		//p = p[chunk:]
	}
	//assign to output

}

/* Final addition in sha256 loop */
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
	h.F.Add2(arrConv(h.out[0]), arrConv(h.chain[0]), arrConv(h.A[64+3]), 32)
	h.F.Add2(arrConv(h.out[1]), arrConv(h.chain[1]), arrConv(h.A[64+2]), 32)
	h.F.Add2(arrConv(h.out[2]), arrConv(h.chain[2]), arrConv(h.A[64+1]), 32)
	h.F.Add2(arrConv(h.out[3]), arrConv(h.chain[3]), arrConv(h.A[64]), 32)
	h.F.Add2(arrConv(h.out[4]), arrConv(h.chain[4]), arrConv(h.E[64+3]), 32)
	h.F.Add2(arrConv(h.out[5]), arrConv(h.chain[5]), arrConv(h.E[64+2]), 32)
	h.F.Add2(arrConv(h.out[6]), arrConv(h.chain[6]), arrConv(h.E[64+1]), 32)
	h.F.Add2(arrConv(h.out[7]), arrConv(h.chain[7]), arrConv(h.E[64]), 32)
}

func (h *Hash) MainLoop(i int) {
	/*
		int sigma0[32], sigma1[32];
		cnf.newVars(sigma0);
		cnf.newVars(sigma1);
		Sigma0(sigma0, A[i+3]);
		Sigma1(sigma1, E[i+3]);
		int f1[32], f2[32];
		cnf.newVars(f1);
		cnf.newVars(f2);
		cnf.ch(f1, E[i+3], E[i+2], E[i+1]);
		cnf.maj3(f2, A[i+3], A[i+2], A[i+1]);
		int T[32];
		cnf.newVars(T);
		cnf.add5(T, E[i], sigma1, f1, k[i], w[i]);
		cnf.add2(E[i+4], A[i], T);
		cnf.add3(A[i+4], T, sigma0, f2);
	*/
	sig0 := h.F.NewVars(32, "sig0")
	sig1 := h.F.NewVars(32, "sig1")
	h.F.Sigma0(sig0, arrConv(h.A[i+3]))
	h.F.Sigma1(sig1, arrConv(h.E[i+3]))
	f1 := h.F.NewVars(32, "f1")
	f2 := h.F.NewVars(32, "f2")
	h.F.If(f1, arrConv(h.E[i+3]), arrConv(h.E[i+2]), arrConv(h.E[i+1]), 32)
	h.F.Maj3(f2, arrConv(h.A[i+3]), arrConv(h.A[i+2]), arrConv(h.A[i+1]), 32)
	t := h.F.NewVars(32, "t")
	h.F.Add5(t, arrConv(h.E[i]), sig1, f1, arrConv(h.K[i]), arrConv(h.W[i]), 32)
	h.F.Add2(arrConv(h.E[i+4]), arrConv(h.A[i]), t, 32)
	h.F.Add3(arrConv(h.A[i+4]), t, sig0, f2, 32)
}

/* Mix p values into w */
func (h *Hash) Mixer1(i int) {
	j := i * 4
	//w[i] = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
	// x := uint32(h.P[j])<<24 | uint32(h.P[j+1])<<16 | uint32(h.P[j+2])<<8 | uint32(h.P[j+3])
	// h.W[i] = [32]int(intToBitVec(x))
	h.W[i] = h.P[j]

}

/* Mix w values */
func (h *Hash) Mixer2(i int) {
	/*
		int s0[32], s1[32];
		cnf.newVars(s0);
		cnf.newVars(s1);
		int r1[32], r2[32];
		cnf.rotr(r1, w[i-15], 7);
		cnf.rotr(r2, w[i-15], 18);
		cnf.xor2(s0+29, r1+29, r2+29, 3);
		cnf.xor3(s0, r1, r2, w[i-15]+3, 29); //shr(w[i-15]+3)
		cnf.rotr(r1, w[i-2], 17);
		cnf.rotr(r2, w[i-2], 19);
		cnf.xor2(s1+22, r1+22, r2+22, 10);
		cnf.xor3(s1, r1, r2, w[i-2]+10, 22); //shr(w[i-2]+10)
		cnf.add4(w[i], w[i-16], s0, w[i-7], s1);
	*/
	s0 := h.F.NewVars(32, "s0")
	s1 := h.F.NewVars(32, "s1")
	r1 := make([]int, 0)
	r2 := make([]int, 0)

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
	// var k = [64][32]int{}
	for i := 0; i < 64; i++ {
		// h.F.NewVars(len(k[i]), "k")
		h.F.NewVars(len(h.K[i]), fmt.Sprintf("k_%d", i))
		// h.F.FixedValue(arrConv(k[i]), int(_K[i]), 32)
		h.F.FixedValue(arrConv(h.K[i]), int(_K[i]), 32)
	}
}

func NewSHA256(input string) Hash {
	h := Hash{}
	h.Rounds = 64
	// h.P = bytesToBitVec([]byte(input))
	h.Input = input
	h.F = NewFormula()
	h.UseIVBlock = true
	// h.in = make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		// h.in[i] = byte(input[i])
	}
	h.InputSize = len(input)*8
	// h.OutputSize = 8
	return h
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
	F          Formula
	Rounds     int    //64
	UseIVBlock bool   //
	InputSize  int    //in bits
	OutputSize int    //
	Input      string //raw input
	Output     string //output string
	chain      [8][32]int //
	out        [8][32]int //output hash value
	P          [][32]int  //input string as bitvec
	//a,b,c,d,e,f,g,h mixing values
	A [70][32]int
	E [70][32]int
	K [64][32]int //k constants
	W [64][32]int //mixing values initialized from input
}
