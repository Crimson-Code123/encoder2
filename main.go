package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {

}

func readDIMACS(name string) Clause {
	file, err := os.OpenFile(name, os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	c := Clause{}
	for x := range strings.Lines(string(data)) {
		_ = x
		//p cnf varcount clausecount

		//
		
	}
	return Clause{}
}

func writeDIMACS(name string, formula Formula) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	s := fmt.Sprintf("p cnf %d %d\n", len(formula.Clauses))
	io.WriteString(file, s)
	
}

type Clause struct {
	Xor bool //not used
	Literals []int
}

type Formula struct {
	Clauses  []Clause
	VarNames []string
}

func (f *Formula)VarCount() int {
	varCount := 0
	for i := 0; i < len(f.Clauses); i++ {
		for j := 0; j < len(f.Clauses[i].Literals); j++ {
			varCount += 1
		}
	}
	return varCount
}
