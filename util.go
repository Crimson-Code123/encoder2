package main

import (
	"crypto/sha256"
	"os"
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

/* Convert fixed length  */
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

/* output is big endian */
func intToBitVec(value uint32) []int {
	x := make([]int, 0)
	for i := 0; i < 32; i++ {
		r := (value >> i) & 1
		x = append(x, int(r))
	}
	return x
}

/*  */
func bitVecToInt(vec []int) uint32 {
	var x int
	for i := 0; i < 32; i++ {
		x = x<<1 | vec[i]
	}
	return uint32(x)
}

/*  */
func shaPaddedMsg(s string) [32]int {
	x := make([]int, 0)
	for i := 0; i < len(s); i++ {
		x = append(x, byteToBitBec(s[i])...)
	}
	return [32]int(x)
}

/* pad sha256 input message */
func padString(r int) string {
	var s string
	s += "€"
	if r != 0 {
		for i := 0; i < r; i++ {
			s += "0"
		}
	}
	return s
}

func byteToBitBec(p byte) []int {
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
		x = append(x, byteToBitBec(p[i])...)
	}
	return x
}

func stringToHex() {

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

func quit() {
	os.Exit(0)
}
