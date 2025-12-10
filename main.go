package main

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

var (
	debug = false
)

func main() {
	// x := readDIMACS("out.dimacs")
	// _ = x
	// fmt.Println(x)
	// writeDIMACS("out2.dimacs", x)
	f := NewFormula()
	
	x := f.NewVars(32, "1")
	r := f.NewVars(32, "2")
	l := f.NewVars(32, "3")
	f.FixedValue(x, 213, 3)
	f.FixedValue(r, 111, 3)
	f.FixedValue(l, 222, 3)
	f.Add2(x, r, l, 32)


	writeDIMACS("test.dimacs", f)
}


