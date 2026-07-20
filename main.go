package main

import (
	"fmt"
	"net"
	"time"
)

var (
	debug  = false
	mixer1 = false
	padding = false
)

func main() {
	// test18()
	// test7()
	// test17()
	// sha256sum("123")
	// test22()
	// return
	var x string
	if !debug && mixer1 {
		fmt.Scanf("%s", x)
		initHash(x, 2)
	} else {
		initHash("123", 1)
	}
	// listener()
	// e := stringToBitVec(sha256sum(("123")))
	// fmt.Println(e)
	// r = fmt.Sprintf("%x", r)
	// fmt.Println(r)
	go server()
	for {
		time.Sleep(time.Second * 10)
	}
}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	// h.Encode(n)
	h.MultiHash()
	writeDIMACS("hash.dimacs", h.F)
}

func server() {
	conn, err := net.Dial("tcp", "127.0.0.1:5656")
	_ = err
	_ = conn
	x := make([]byte, 1024)
	n, err := conn.Read(x)
	fmt.Printf("N: %d | X: %s\n", n, x)

}




