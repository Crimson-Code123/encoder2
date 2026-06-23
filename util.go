package main

import (
	"fmt"
	sha256 "main/reference"
	"os"
	"strconv"
)

/*
	unsigned long long ml = size * 8;
	int i = size;
	vector<unsigned char> M(m, m+size);
	M.push_back(0x80);
	while( M.size() % 64 != 56 ) M.push_back(0x00);
	//big endian
	M.push_back((ml >> 56) & 0xFF);
	M.push_back((ml >> 48) & 0xFF);
	M.push_back((ml >> 40) & 0xFF);
	M.push_back((ml >> 32) & 0xFF);
	M.push_back((ml >> 24) & 0xFF);
	M.push_back((ml >> 16) & 0xFF);
	M.push_back((ml >> 8) & 0xFF);
	M.push_back((ml) & 0xFF);
	size = M.size();
	unsigned w[64]; //
	Process the message in successive 512-bit chunks:
	   break message into 512-bit chunks
	for( int chunk=0; chunk<size; chunk+=64 )
	{
		for( int i=0; i<16; i++ ) {
			w[i] = 0;
			for( int j=0; j<4; j++ ) {
				w[i] = (w[i] << 8) | M[chunk + i*4 + j];
			}
		}
		sha256_comp(w, H, rounds, true);
	}
*/



/* Util functions */

func assertOutput() bool {

	return false
}

func sha256sum(data string) string {
	h := sha256.Sum256([]byte(data))
	x := make([]byte, 0)
	for i := 0; i < len(h); i++ {
		x = append(x, h[i])
	}
	return fmt.Sprintf("%x", string(x))
}

func shaPaddedMsg(s string) [32]int {
	x := make([]int, 0)
	for i := 0; i < len(s); i++ {
		x = append(x, byteToBitVec(s[i])...)
	}
	return [32]int(x)
}

/* pad sha256 input message */
func padString(x string) string {
	var s string
	s += "€" //0x80
	// while( M.size() % 64 != 56 ) M.push_back(0x00);
	// fmt.Println(len(s))
	// for len(s)%64 != 56 {
	for len(s)%16 != 15 {
		s += "0"
	}
	// if r != 0 {
	// 	for i := 0; i < r; i++ {
	// 		s += "0"
	// 	}
	// }
	//pad size, 64 bit big endian
	var z uint64
	z = uint64(len(x) + len(s))
	zb := bitVecToString(reverse(int64ToBitVec(z)))
	x += s
	x += zb
	return x
}

/* get difficulty bits, todo: online updates. */
func calcDifficulty(x string) int {
	/* todo: fix output bits (get from btc difficulty)
	64 chars, 256 bits (152 bits as of may 2026)
	0000000000000000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF
	*/
	diff := len(x) * 8
	return diff
}

/* Convert to fixed length  */
func arrConv(item [32]int) []int {
	r := make([]int, 0)
	for i := 0; i < len(item); i++ {
		r = append(r, item[i])
	}
	return r
}

func arrConv8(item [8]int) []int {
	r := make([]int, 0)
	for i := 0; i < len(item); i++ {
		r = append(r, item[i])
	}
	return r
}

/* output is big endian */
func intToBitVec(value uint32) []int {
	x := make([]int, 0)
	for i := 0; i < 32; i++ {
		r := (value >> i) & 1
		x = append(x, int(r))
	}
	return x
}

/* output is big endian */
func int64ToBitVec(value uint64) []int {
	x := make([]int, 0)
	for i := 0; i < 64; i++ {
		r := (value >> i) & 1
		x = append(x, int(r))
	}
	return x
}

/* works */
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

/* works */
func bitVecToByte(a []int) byte {
	var x byte
	a = reverse(a)
	for i := 0; i < 8; i++ {
		x = x<<1 | byte(a[i])
	}
	return x
}

/* works */
func byteToBitVec(p byte) []int {
	x := make([]int, 0)
	for i := 0; i < 8; i++ {
		r := (p >> i) & 1
		x = append(x, int(r))
	}
	return x
}

/* convert byte string to bit vec */
func bytesToBitVec(p []byte) []int {
	x := make([]int, 0)
	for i := 0; i < len(p); i++ {
		x = append(x, byteToBitVec(p[i])...)
	}
	return x
}

/* convert input to continuous array */
func hashOutput(x [8][32]int) []int {
	r := make([]int, 0)
	for i := 0; i < len(x); i++ {
		s := arrConv(x[i])
		for j := 0; j < len(s); j++ {
			r = append(r, s[j])
		}
	}
	return r
}

/* works */
func bitVecToBytes(a []int) []byte {
	r := make([]byte, 0)
	for i := 0; i < len(a); i += 8 {
		r = append(r, bitVecToByte(a[i:i+8]))
	}
	return r
}

func stringToBitVec(s string) []int {
	return bytesToBitVec([]byte(s))
}

func bitVecToString(a []int) string {
	return string(bitVecToBytes(a))
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

func itoa(i int) string {
	return strconv.Itoa(i)
}

func quit() {
	os.Exit(0)
}

// func _8to32(item... [8]int) [32]int {
// 	r := make([]int, 0)
// 	if len(item) != 4 {0000000
// 		return [32]int(r)
// 	}
// 	for i := 0; i < len(item); i++ {
// 		r = append(r, item[i])
// 	}
// 	return [32]int(r)
// }

/*  */
// func intToBytes(i int) []byte {
// 	r := make([]byte, 0)
// 	for x := 0; x < 4; x++ {
// 		var t byte
// 		t = t << 1
// 		r = append(r, t)
// 	}
// 	return r
// }
