// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SHA256 block step.

package sha256

import (
	"fmt"
	"math/bits"
)

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

func blockGeneric(dig *digest, p []byte) {
	fmt.Printf("")
	var w [64]uint32
	//assign IV
	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
	for len(p) >= chunk {
		for i := 0; i < 16; i++ {
			j := i * 4
			var x uint32
			aa := uint32(p[j])
			bb := uint32(p[j+1])
			cc := uint32(p[j+2])
			dd := uint32(p[j+3])
			fmt.Printf("Non aa: %d | bb: %d | cc: %d | dd: %d\n", aa, bb, cc, dd)
			aaa := reverse(intToBitVec(aa))
			bbb := reverse(intToBitVec(bb))
			ccc := reverse(intToBitVec(cc))
			ddd := reverse(intToBitVec(dd))
			if aa != 0 {
				fmt.Println(aaa, bbb, ccc, ddd)
			}
			zz := uint32(p[j]) << 24
			yy := uint32(p[j+1]) << 16
			xx := uint32(p[j+2]) << 8
			ww := uint32(p[j+3])
			// fmt.Printf("Shifted zz: %d | yy: %d | xx: %d | ww: %d\n", zz, yy, xx, ww)
			// x = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
			x = zz | yy | xx | ww
			// fmt.Println("x:", x)
			zzz := reverse(intToBitVec(zz))
			yyy := reverse(intToBitVec(yy))
			xxx := reverse(intToBitVec(xx))
			www := reverse(intToBitVec(ww))
			if aa != 0 {
				fmt.Println(zzz, yyy, xxx, www)
			}
			// fmt.Printf("\n")
			// fmt.Println("Final:", intToBitVec(x))
			w[i] = x
		}
		for i := 16; i < 64; i++ {
			v1 := w[i-2]
			t1 := (bits.RotateLeft32(v1, -17)) ^ (bits.RotateLeft32(v1, -19)) ^ (v1 >> 10)
			v2 := w[i-15]
			t2 := (bits.RotateLeft32(v2, -7)) ^ (bits.RotateLeft32(v2, -18)) ^ (v2 >> 3)
			w[i] = t1 + w[i-7] + t2 + w[i-16]
		}
		a, b, c, d, e, f, g, h := h0, h1, h2, h3, h4, h5, h6, h7
		for i := 0; i < 64; i++ {
			t1 := h + ((bits.RotateLeft32(e, -6)) ^ (bits.RotateLeft32(e, -11)) ^ (bits.RotateLeft32(e, -25))) + ((e & f) ^ (^e & g)) + _K[i] + w[i]
			t2 := ((bits.RotateLeft32(a, -2)) ^ (bits.RotateLeft32(a, -13)) ^ (bits.RotateLeft32(a, -22))) + ((a & b) ^ (a & c) ^ (b & c))
			h = g
			g = f
			f = e
			e = d + t1
			d = c
			c = b
			b = a
			a = t1 + t2
			// fmt.Println(a, b, c, d, e, f, g, h)
		}
		//final addition
		h0 += a
		h1 += b
		h2 += c
		h3 += d
		h4 += e
		h5 += f
		h6 += g
		h7 += h
		// fmt.Println(h0, h1, h2, h3, h4, h5, h6, h7)
		p = p[chunk:]
	}
	//assign
	dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7] = h0, h1, h2, h3, h4, h5, h6, h7
}

func intToBitVec(value uint32) []int {
	x := make([]int, 0)
	for i := 0; i < 32; i++ {
		r := (value >> i) & 1
		x = append(x, int(r))
	}
	return x
}

func bitVecToInt(vec []int) uint32 {
	var x int
	if len(vec) > 32 {
		return 0
	}
	for i := 0; i < len(vec); i++ {
		x = x<<1 | vec[i]
	}
	return uint32(x)
}

func reverse(a []int) []int {
	x := make([]int, len(a))
	r := len(a) - 1
	for i := 0; i < len(a); i++ {
		x[r] = a[i]
		r -= 1
	}
	return x
}
