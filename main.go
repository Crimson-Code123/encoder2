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
			fmt.Printf("Read file with %s Vars and %s Clauses\n", y[2], y[3])
			i += 1
		} else {
			if y[0] == "c" { //c varname varid
				r, _ := strconv.Atoi(y[2])
				f.VarNames[y[1]] = r
			} else {
				c := Clause{}
				for z := 0; z < len(y); z++ {
					r, _ := strconv.Atoi(y[z])
					c.Literals = append(c.Literals, r)
				}
				f.Clauses = append(f.Clauses, c)
			}

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
	s := fmt.Sprintf("p cnf %d %d\n", formula.VarCount(), len(formula.Clauses))
	io.WriteString(file, s)
	for i := 0; i < len(formula.Clauses); i++ {
		for j := 0; j < len(formula.Clauses[i].Literals); j++ {
			s := fmt.Sprintf("%d ", formula.Clauses[i].Literals[j])
			io.WriteString(file, s)
		}
		io.WriteString(file, "0\n")
	}
	for x, y := range formula.VarNames {
		s := fmt.Sprintf("c %s %d\n", x, y)
		io.WriteString(file, s)
	}
}

type Clause struct {
	Literals []int
}

type Formula struct {
	Clauses  []Clause
	VarNames map[string]int
}

func NewFormula() Formula {
	f := Formula{}
	f.VarNames = make(map[string]int)
	return f
}

func (f *Formula) NewVars() {

}

func (f *Formula) AddClause() {

}

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
	return varCount
}
