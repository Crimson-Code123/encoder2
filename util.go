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
