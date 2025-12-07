package main

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

var (
	debug = false
)

func main() {
	x := readDIMACS("out.dimacs")
	_ = x
	// fmt.Println(x)
	// writeDIMACS("out2.dimacs", x)
	
}


