package main

import (
	"encoding/binary"
	"fmt"
	"math/bits"
)

var w [64]uint32

func padded(p string) string {
	return Sum256([]byte(p))
}

func bgh(h *digest, p []byte) {
	for len(p) >= chunk {
		for i := 0; i < 16; i++ {
			j := i * 4
			var x uint32
			zz := uint32(p[j]) << 24
			yy := uint32(p[j+1]) << 16
			xx := uint32(p[j+2]) << 8
			ww := uint32(p[j+3])
			x = zz | yy | xx | ww
			w[i] = x
			fmt.Printf("W[%d]: %d\n", i, x)
		}
		// for i := 16; i < 64; i++ {
		// 	v1 := w[i-2]
		// 	t1 := (bits.RotateLeft32(v1, -17)) ^ (bits.RotateLeft32(v1, -19)) ^ (v1 >> 10)
		// 	v2 := w[i-15]
		// 	t2 := (bits.RotateLeft32(v2, -7)) ^ (bits.RotateLeft32(v2, -18)) ^ (v2 >> 3)
		// 	w[i] = t1 + w[i-7] + t2 + w[i-16]
		// }
		p = p[chunk:]
	}
}

func blockGeneric(dig *digest, p []byte) {
	var w [64]uint32
	//assign IV
	h0, h1, h2, h3, h4, h5, h6, h7 := dig.h[0], dig.h[1], dig.h[2], dig.h[3], dig.h[4], dig.h[5], dig.h[6], dig.h[7]
	for len(p) >= chunk {
		for i := 0; i < 16; i++ {
			j := i * 4
			var x uint32
			// aa := uint32(p[j])
			// bb := uint32(p[j+1])
			// cc := uint32(p[j+2])
			// dd := uint32(p[j+3])
			// fmt.Printf("Non aa: %d | bb: %d | cc: %d | dd: %d\n", aa, bb, cc, dd)
			// aaa := reverse(intToBitVec(aa))
			// bbb := reverse(intToBitVec(bb))
			// ccc := reverse(intToBitVec(cc))
			// ddd := reverse(intToBitVec(dd))
			// if aa != 0 {
			// 	fmt.Println(aaa, bbb, ccc, ddd)
			// }
			zz := uint32(p[j]) << 24
			yy := uint32(p[j+1]) << 16
			xx := uint32(p[j+2]) << 8
			ww := uint32(p[j+3])
			// fmt.Printf("Shifted zz: %d | yy: %d | xx: %d | ww: %d\n", zz, yy, xx, ww)
			// x = uint32(p[j])<<24 | uint32(p[j+1])<<16 | uint32(p[j+2])<<8 | uint32(p[j+3])
			x = zz | yy | xx | ww
			// fmt.Println("x:", x)
			// zzz := reverse(intToBitVec(zz))
			// yyy := reverse(intToBitVec(yy))
			// xxx := reverse(intToBitVec(xx))
			// www := reverse(intToBitVec(ww))
			// if aa != 0 {
			// 	fmt.Println(zzz, yyy, xxx, www)
			// }
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

func (d *digest) Write(p []byte) (nn int, err error) {
	nn = len(p)
	d.len += uint64(nn)
	if d.nx > 0 {
		n := copy(d.x[d.nx:], p)
		// fmt.Println("dx: ", len(d.x[d.nx:]), d.x[d.nx:])
		d.nx += n
		if d.nx == chunk {
			fmt.Println("d.x: ", len(d.x[:]), d.x[:])
			bgh(d, d.x[:])
			d.nx = 0
		} else {

		}
		p = p[n:]
	}
	// if len(p) != 0 {
	// 	fmt.Println("zx:", len(p), p)
	// }
	if len(p) >= chunk {
		n := len(p) &^ (chunk - 1)
		fmt.Println("n: ", n, p[:n]) //d.x[:],
		bgh(d, p[:n])
		p = p[n:]
		// fmt.Println("z:", len(p), p)
	}
	if len(p) > 0 { //
		d.nx = copy(d.x[:], p)
	}
	return
}

func (d *digest) checkSum() {
	length := d.len
	// Padding. Add a 1 bit and 0 bits until 56 bytes mod 64.
	var tmp [64 + 8]byte // padding + length buffer
	tmp[0] = 0x80
	var t uint64
	if length%64 < 56 {
		t = 56 - length%64
	} else {
		t = 64 + 56 - length%64
	}
	// Length in bits.
	length <<= 3
	padlen := tmp[:t+8]
	binary.BigEndian.PutUint64(padlen[t+0:], length)
	// fmt.Println("Padlen:", len(padlen), padlen)
	//
	d.Write(padlen)

	// if d.nx != 0 {
	// 	panic("d.nx != 0")
	// }
	// var digest [Size]byte
	// binary.BigEndian.PutUint32(digest[0:], d.h[0])
	// binary.BigEndian.PutUint32(digest[4:], d.h[1])
	// binary.BigEndian.PutUint32(digest[8:], d.h[2])
	// binary.BigEndian.PutUint32(digest[12:], d.h[3])
	// binary.BigEndian.PutUint32(digest[16:], d.h[4])
	// binary.BigEndian.PutUint32(digest[20:], d.h[5])
	// binary.BigEndian.PutUint32(digest[24:], d.h[6])
	// if !d.is224 {
	// 	binary.BigEndian.PutUint32(digest[28:], d.h[7])
	// }
	// fmt.Println(digest)
	// return digest

}

// Sum256 returns the SHA256 checksum of the data.
func Sum256(data []byte) string {
	var d digest
	d.Reset()
	d.Write(data)
	d.checkSum()
	return ""
	// return string(d.x)
	// return d.checkSum()
}

// digest represents the partial evaluation of a checksum.
type digest struct {
	h     [8]uint32
	x     [chunk]byte
	nx    int
	len   uint64
	is224 bool // mark if this digest is SHA-224
}

// The size of a SHA256 checksum in bytes.
const Size = 32

// The size of a SHA224 checksum in bytes.
const Size224 = 28

// The blocksize of SHA256 and SHA224 in bytes.
const BlockSize = 64

const (
	chunk = 64
	init0 = 0x6A09E667
	init1 = 0xBB67AE85
	init2 = 0x3C6EF372
	init3 = 0xA54FF53A
	init4 = 0x510E527F
	init5 = 0x9B05688C
	init6 = 0x1F83D9AB
	init7 = 0x5BE0CD19
)

func (d *digest) Reset() {
	if !d.is224 {
		d.h[0] = init0
		d.h[1] = init1
		d.h[2] = init2
		d.h[3] = init3
		d.h[4] = init4
		d.h[5] = init5
		d.h[6] = init6
		d.h[7] = init7
	}
	d.nx = 0
	d.len = 0
}
