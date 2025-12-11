package main

import (
	"crypto/sha256"
)

/* Util functions */

func sha256sum(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return string(h.Sum([]byte(data)))
}

/* Convert fixed length  */
func arrayConv(item [32]int) []int {
	r := make([]int, 0)
	for i := 0; i < len(item); i++ {
		r = append(r, item[i])
	}
	return r
}


func intToBitVec() []int {
	return make([]int, 0)
}
