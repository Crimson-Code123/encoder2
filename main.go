package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f := NewFormula()
	f.AddClause(NewClause([]int{1, -2, -3, 4}))
	f.AddClause(NewClause([]int{55, -33, 22, 11}))
	f.VarNames["123"] = 2
	// writeDIMACS("out.dimacs", f)
	x := readDIMACS("out.dimacs")
	fmt.Println(x)
}

func readDIMACS(name string) Formula {
	file, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	f := NewFormula()
	i := 0
	for x := range strings.Lines(string(data)) {
		y := strings.Split(x, " ")
		if i == 0 { //p cnf varcount clausecount
			fmt.Printf("Read file with %s Vars and %s Clauses\n", y[2], strings.Split(y[3], "\n")[0])
			i += 1
		} else {
			if y[0] == "c" { //c varname varid
				r, err := strconv.Atoi(strings.Split(y[2], "\n")[0])
				if err != nil {
					fmt.Println("Named var: ", err)
				}
				f.VarNames[y[1]] = r
			} else {
				c := Clause{}
				for z := 0; z < len(y)-1; z++ {
					r, err := strconv.Atoi(y[z])
					if err != nil {
						fmt.Println(err)
					}
					c.Literals = append(c.Literals, r)
				}
				f.Clauses = append(f.Clauses, c)
			}
			i += 1
		}
	}
	return f
}

func writeDIMACS(name string, formula Formula) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	s := fmt.Sprintf("p cnf %d %d \n", formula.VarCount(), len(formula.Clauses))
	io.WriteString(file, s)
	for i := 0; i < len(formula.Clauses); i++ {
		for j := 0; j < len(formula.Clauses[i].Literals); j++ {
			s := fmt.Sprintf("%d ", formula.Clauses[i].Literals[j])
			io.WriteString(file, s)
		}
		io.WriteString(file, "0\n")
	}
	for x, y := range formula.VarNames {
		s := fmt.Sprintf("c %s %d \n", x, y)
		io.WriteString(file, s)
	}
}

/*
 AdderType: Ripple carry
 MultiAdderType: Two operand
 PB Method: Sequential counter
*/

type Formula struct {
	Clauses  []Clause
	VarNames map[string]int
	VarID    int
	varCount int
}

type Clause struct {
	Literals []int
}

func NewFormula() Formula {
	f := Formula{}
	f.VarNames = make(map[string]int)
	return f
}

func NewClause(vars []int) Clause {
	c := Clause{}
	c.Literals = vars
	return c
}

func (f *Formula) Add2(z []int, x []int, y []int, n int) {
	r := n - 1
	c := make([]int, r)
	f.NewVars(c, r, "")

	f.HalfAdder(c, z, x, y, 1)

	// f.FullAdder(c+1, z+1, x+1, y+1, c, n-2)

	// f.Xor3(z+n-1, x+n-1, y+n-1, c+n-2, 1)

}

func (f *Formula) Add3() {

}

func (f *Formula) Add4() {

}

func (f *Formula) Add5() {

}

func (f *Formula) HalfAdder(c []int, s []int, x []int, y []int, n int) {
	f.Xor2(s, x, y, n)
	f.And2(c, x, y, n)
}

func (f *Formula) FullAdder(c []int, s []int, x []int, y []int, t []int, n int) {
	f.Xor3(s, x, y, t, n)
	f.Maj3(c, x, y, t, n)
}

func (f *Formula)Or2(z []int, x []int, y []int, n int) {

}

func (f *Formula) And2(z []int, x []int, y []int, n int) {

}

func (f *Formula) Maj3(z []int, x []int, y []int, t []int, n int) {
	for i := 0; i < n; i++ {
		f.AddClause(NewClause([]int{-z[i], x[i], y[i]}))
		f.AddClause(NewClause([]int{-z[i], x[i], t[i]}))
		f.AddClause(NewClause([]int{-z[i], y[i], t[i]}))
		f.AddClause(NewClause([]int{z[i], -y[i], -t[i]}))
		f.AddClause(NewClause([]int{z[i], -x[i], -t[i]}))
		f.AddClause(NewClause([]int{z[i], -x[i], -y[i]}))
	}
}

func (f *Formula) Xor2(z []int, x []int, y []int, n int) {
	for i := 0; i < n; i++ {
		f.AddClause(NewClause([]int{-z[i], -x[i]}))
	}
}

func (f *Formula) Xor3(z []int, x []int, y []int, t []int, n int) {
	for i := 0; i < n; i++ {

	}
}

func (f *Formula) NewVars(ivec []int, vars int, name string) {
	for i := 0; i < vars; i++ {
		f.VarID += 1
		ivec[i] = f.VarID
	}
	if name != "" {
		f.VarNames[name] = ivec[0]
	}
	f.varCount += vars
}

func (f *Formula) AddClause(c Clause) {
	f.Clauses = append(f.Clauses, c)
}

// func (f *Formula)AddClause(c []int) {
// 	x := Clause{}
// 	x.Literals = c
// 	f.Clauses = append(f.Clauses, )
// }

func (f *Formula) FixedValue() {

}

func (f *Formula) AppendFormula(fx Formula) {

}

func (f *Formula) VarCount() int {
	varCount := 0
	for i := 0; i < len(f.Clauses); i++ {
		for j := 0; j < len(f.Clauses[i].Literals); j++ {
			varCount += 1
		}
	}
	f.varCount = varCount //
	return varCount
}
