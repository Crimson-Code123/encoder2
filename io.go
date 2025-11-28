package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/* IO read/write for DIMACS.  */

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
