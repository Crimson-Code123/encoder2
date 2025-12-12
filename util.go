package main

import (
	"crypto/sha256"
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

/* output is little endian */
func intToBitVec(value uint32) []int {
	x := make([]int, 0)
	for i := 0; i < 32; i++ {
		r := (value >> i) & 1
		x = append(x, int(r))
	}
	return x
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
