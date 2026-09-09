package main

import (
	"fmt"
	"net"
)

/*

*/

var (
	debug   = false
	mixer1  = false
	padding = false
	r = "12345678"
	// r = "1234567812345678123456781234567812345678123456781234567812345678"
)

func main() {
	// test18()
	// test7()
	// test17()
	// sha256sum("123")
	// test22()
	// e := "12376575785656787777657565675134134375555555555555555555555555555555555555555565555675566578667757687687"
	// Sum256([]byte(e))
	// fmt.Println(sha256sum(e))
	// return
	var x string
	if !debug && mixer1 {
		fmt.Scanln("%s", x)
		initHash(x, 2)
	} else {
		initHash(r, 1)
	}
	
	// listener()
	// e := stringToBitVec(sha256sum(("123")))
	// fmt.Println(e)
	// r = fmt.Sprintf("%x", r)
	// fmt.Println(r)
	// go server()
	// for {
	// 	time.Sleep(time.Second * 10)
	// }
}

func initHash(input string, n int) {
	h := NewSHA256(input, n)
	// h.Encode(n)
	h.MultiHash()
	transform(h.F)
	writeDIMACS("hash.dimacs", h.F)
}

func server() {
	conn, err := net.Dial("tcp", "127.0.0.1:5656")
	_ = err
	_ = conn
	x := make([]byte, 1024)
	_ = x
	// x := new([]byte
	n, err := conn.Read(x)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("N: %d | X: %s\n", n, x)
}




