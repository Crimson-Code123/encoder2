package main

import "fmt"

/* Assert functionality of CNF functions */
func testSampleFormula() {
	fmt.Println()
}

/* test io functionality */
func test1() {
	x := readDIMACS("out.dimacs")
	fmt.Println(x)
	writeDIMACS("out2.dimacs", x)
}

/* test cnf functionality */
func test2() {
	f := NewFormula()
	x := f.NewVars(32, "1")
	r := f.NewVars(32, "2")
	l := f.NewVars(32, "3")
	f.FixedValue(x, 213, 32)
	f.FixedValue(r, 111, 32)
	f.FixedValue(l, 222, 32)
	f.Add2(x, r, l, 32)
	f.Xor2(x, r, l, 32)
	f.Sigma0(x, r)
	fmt.Println(x, r)
	x = f.RotateR(x, r, 2, 32)
	fmt.Println(x, r)
	writeDIMACS("test.dimacs", f)
}

