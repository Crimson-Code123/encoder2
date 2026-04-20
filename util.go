package main

import (
	sha256 "main/reference"
	"os"
	"strconv"
)

/* Util functions */

func sha256sum(data string) string {
	h := sha256.Sum256([]byte(data))
	x := make([]byte, 0)
	for i := 0; i < len(h); i++ {
		x = append(x, h[i])
	}
	return string(x)
}


/*  */
func shaPaddedMsg(s string) [32]int {
	x := make([]int, 0)
	for i := 0; i < len(s); i++ {
		x = append(x, byteToBitVec(s[i])...)
	}
	return [32]int(x)
}

/* pad sha256 input message */
func padString(r int) string {
	var s string
	s += "€" //0x80
	if r != 0 {
		for i := 0; i < r; i++ {
			s += "0" 
		}
	}
	return s
}


// func _8to32(item... [8]int) [32]int {
// 	r := make([]int, 0)
// 	if len(item) != 4 {
// 		return [32]int(r)
// 	}
// 	for i := 0; i < len(item); i++ {
// 		r = append(r, item[i])
// 	}
// 	return [32]int(r)
// }


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

/* works */
func bitVecToInt(vec []int) uint32 {
	var x int
	for i := 0; i < 32; i++ {
		x = x<<1 | vec[i]
	}
	return uint32(x)
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

/* works */
func bitVecToByte(a []int) byte {
	var x byte
	a = reverse(a)
	for i := 0; i < 8; i++ {
		x = x<<1 | byte(a[i])
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

/*  */
func intToBytes(i int) []byte {
	r := make([]byte, 0)
	for x := 0; x < 4; x++ {
		var t byte
		t = t << 1 
		r = append(r, t)
	}
	return r
}

func bitVecToBytes(a []int) []byte {
	r := make([]byte, 0)
	for i := 0; i < len(a); i++ {
		r = append(r, intToBytes(i)...)
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
