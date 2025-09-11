package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// 47b60685dfd5ceed368238b9611a105a46f80d7973e72a29b87ec8542f22fdf5

func main() {
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
	Vars     map[int][]int
	VarID    int
	varCount int
}

type Clause struct {
	Literals []int
}

type Var struct {
	Name    string
	Value   uint
	Length  int //32 bit
	ID      int
	BoolVec []bool
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

	f.FullAdder(c[1:], z[1:], x[1:], y[1:], c, n-2)

	f.Xor3(z, x, y, c, 1) //

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

func (f *Formula) Or2(z []int, x []int, y []int, n int) {

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
		f.AddClause(NewClause([]int{-z[i], -x[i], -y[i]}))
		f.AddClause(NewClause([]int{z[i], -x[i], y[i]}))
		f.AddClause(NewClause([]int{z[i], x[i], -y[i]}))
		f.AddClause(NewClause([]int{-z[i], x[i], y[i]}))
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

func (f *Formula)AddClauseVec(c []int) {
	x := Clause{}
	x.Literals = c
	f.Clauses = append(f.Clauses, x)
}

func (f *Formula) FixedValue(z []int, value int, n int) {
	for i := 0; i < n; i++ {
		var x int
		if z[i] < 0 {
			x = (value >> i) & 1
		} else {
			x = (value >> i) & 1
		}
		fmt.Println(x, -x)
		f.AddClause(NewClause([]int{int(x)}))
	}
}

func (f *Formula) AppendFormula(fx Formula) {

}

func (f *Formula) VarCount() int {
	varCount := 0
	for i := 0; i < len(f.Clauses); i++ {
		for j := 0; j < len(f.Clauses[i].Literals); j++ {
			if f.Clauses[i].Literals[j] > varCount {
				varCount = f.Clauses[i].Literals[j]
			}
		}
	}
	f.varCount = varCount //
	return varCount
}

type Hash struct {
	Rounds     int
	UseIVBlock bool
	InputSize  int
	OutputSize int
	Form       Formula
	w          []byte
	chain      []byte
	out        []byte
}

/*
value: binary encoded into sat variables
x/z: pointer address
z[i]/x[i]: variable id

Newvar: varid: 17034 | x: 17013 | x[i]: 17034
Newvar: varid: 17035 | x: 17013 | x[i]: 17035
Newvar: varid: 17036 | x: 17013 | x[i]: 17036
Newvar: varid: 17037 | x: 17013 | x[i]: 17037
Newvar: varid: 17038 | x: 17013 | x[i]: 17038
Newvar: varid: 17039 | x: 17013 | x[i]: 17039
Newvar: varid: 17040 | x: 17013 | x[i]: 17040
Newvar: varid: 17041 | x: 17013 | x[i]: 17041
Newvar: varid: 17042 | x: 17013 | x[i]: 17042
Newvar: varid: 17043 | x: 17013 | x[i]: 17043
Defaults: z: 1985 | x: 16887 | y: 16950 | c: 17013
full: z: 1986 | x: 16888 | y: 16951 | c: 17013
xor3: z: 2016 | x: 16918 | y: 16981 | c: 17045
Expansion loop 63
Newvar: varid: 17044 | x: 17044 | x[i]: 17044
Newvar: varid: 17045 | x: 17044 | x[i]: 17045
Newvar: varid: 17046 | x: 17044 | x[i]: 17046
Newvar: varid: 17047 | x: 17044 | x[i]: 17047
Newvar: varid: 17048 | x: 17044 | x[i]: 17048
Newvar: varid: 17049 | x: 17044 | x[i]: 17049
Newvar: varid: 17050 | x: 17044 | x[i]: 17050
Newvar: varid: 17051 | x: 17044 | x[i]: 17051
Newvar: varid: 17052 | x: 17044 | x[i]: 17052
Newvar: varid: 17053 | x: 17044 | x[i]: 17053
Newvar: varid: 17054 | x: 17044 | x[i]: 17054
Newvar: varid: 17055 | x: 17044 | x[i]: 17055
Newvar: varid: 17056 | x: 17044 | x[i]: 17056
Newvar: varid: 17057 | x: 17044 | x[i]: 17057
Newvar: varid: 17058 | x: 17044 | x[i]: 17058
Newvar: varid: 17059 | x: 17044 | x[i]: 17059
Newvar: varid: 17060 | x: 17044 | x[i]: 17060
Newvar: varid: 17061 | x: 17044 | x[i]: 17061
Newvar: varid: 17062 | x: 17044 | x[i]: 17062
Newvar: varid: 17063 | x: 17044 | x[i]: 17063
Newvar: varid: 17064 | x: 17044 | x[i]: 17064
Newvar: varid: 17065 | x: 17044 | x[i]: 17065
Newvar: varid: 17066 | x: 17044 | x[i]: 17066
Newvar: varid: 17067 | x: 17044 | x[i]: 17067
Newvar: varid: 17068 | x: 17044 | x[i]: 17068
Newvar: varid: 17069 | x: 17044 | x[i]: 17069
Newvar: varid: 17070 | x: 17044 | x[i]: 17070
Newvar: varid: 17071 | x: 17044 | x[i]: 17071
Newvar: varid: 17072 | x: 17044 | x[i]: 17072
Newvar: varid: 17073 | x: 17044 | x[i]: 17073
Newvar: varid: 17074 | x: 17044 | x[i]: 17074
Newvar: varid: 17075 | x: 17044 | x[i]: 17075
Newvar: varid: 17076 | x: 17076 | x[i]: 17076
Newvar: varid: 17077 | x: 17076 | x[i]: 17077
Newvar: varid: 17078 | x: 17076 | x[i]: 17078
Newvar: varid: 17079 | x: 17076 | x[i]: 17079
Newvar: varid: 17080 | x: 17076 | x[i]: 17080
Newvar: varid: 17081 | x: 17076 | x[i]: 17081
Newvar: varid: 17082 | x: 17076 | x[i]: 17082
Newvar: varid: 17083 | x: 17076 | x[i]: 17083
Newvar: varid: 17084 | x: 17076 | x[i]: 17084
Newvar: varid: 17085 | x: 17076 | x[i]: 17085
Newvar: varid: 17086 | x: 17076 | x[i]: 17086
Newvar: varid: 17087 | x: 17076 | x[i]: 17087
Newvar: varid: 17088 | x: 17076 | x[i]: 17088
Newvar: varid: 17089 | x: 17076 | x[i]: 17089
Newvar: varid: 17090 | x: 17076 | x[i]: 17090
Newvar: varid: 17091 | x: 17076 | x[i]: 17091
Newvar: varid: 17092 | x: 17076 | x[i]: 17092
Newvar: varid: 17093 | x: 17076 | x[i]: 17093
Newvar: varid: 17094 | x: 17076 | x[i]: 17094
Newvar: varid: 17095 | x: 17076 | x[i]: 17095
Newvar: varid: 17096 | x: 17076 | x[i]: 17096
Newvar: varid: 17097 | x: 17076 | x[i]: 17097
Newvar: varid: 17098 | x: 17076 | x[i]: 17098
Newvar: varid: 17099 | x: 17076 | x[i]: 17099
Newvar: varid: 17100 | x: 17076 | x[i]: 17100
Newvar: varid: 17101 | x: 17076 | x[i]: 17101
Newvar: varid: 17102 | x: 17076 | x[i]: 17102
Newvar: varid: 17103 | x: 17076 | x[i]: 17103
Newvar: varid: 17104 | x: 17076 | x[i]: 17104
Newvar: varid: 17105 | x: 17076 | x[i]: 17105
Newvar: varid: 17106 | x: 17076 | x[i]: 17106
Newvar: varid: 17107 | x: 17076 | x[i]: 17107
Newvar: varid: 17108 | x: 17108 | x[i]: 17108
Newvar: varid: 17109 | x: 17108 | x[i]: 17109
Newvar: varid: 17110 | x: 17108 | x[i]: 17110
Newvar: varid: 17111 | x: 17108 | x[i]: 17111
Newvar: varid: 17112 | x: 17108 | x[i]: 17112
Newvar: varid: 17113 | x: 17108 | x[i]: 17113
Newvar: varid: 17114 | x: 17108 | x[i]: 17114
Newvar: varid: 17115 | x: 17108 | x[i]: 17115
Newvar: varid: 17116 | x: 17108 | x[i]: 17116
Newvar: varid: 17117 | x: 17108 | x[i]: 17117
Newvar: varid: 17118 | x: 17108 | x[i]: 17118
Newvar: varid: 17119 | x: 17108 | x[i]: 17119
Newvar: varid: 17120 | x: 17108 | x[i]: 17120
Newvar: varid: 17121 | x: 17108 | x[i]: 17121
Newvar: varid: 17122 | x: 17108 | x[i]: 17122
Newvar: varid: 17123 | x: 17108 | x[i]: 17123
Newvar: varid: 17124 | x: 17108 | x[i]: 17124
Newvar: varid: 17125 | x: 17108 | x[i]: 17125
Newvar: varid: 17126 | x: 17108 | x[i]: 17126
Newvar: varid: 17127 | x: 17108 | x[i]: 17127
Newvar: varid: 17128 | x: 17108 | x[i]: 17128
Newvar: varid: 17129 | x: 17108 | x[i]: 17129
Newvar: varid: 17130 | x: 17108 | x[i]: 17130
Newvar: varid: 17131 | x: 17108 | x[i]: 17131
Newvar: varid: 17132 | x: 17108 | x[i]: 17132
Newvar: varid: 17133 | x: 17108 | x[i]: 17133
Newvar: varid: 17134 | x: 17108 | x[i]: 17134
Newvar: varid: 17135 | x: 17108 | x[i]: 17135
Newvar: varid: 17136 | x: 17108 | x[i]: 17136
Newvar: varid: 17137 | x: 17108 | x[i]: 17137
Newvar: varid: 17138 | x: 17108 | x[i]: 17138
Newvar: varid: 17139 | x: 17108 | x[i]: 17139
Newvar: varid: 17140 | x: 17140 | x[i]: 17140
Newvar: varid: 17141 | x: 17140 | x[i]: 17141
Newvar: varid: 17142 | x: 17140 | x[i]: 17142
Newvar: varid: 17143 | x: 17140 | x[i]: 17143
Newvar: varid: 17144 | x: 17140 | x[i]: 17144
Newvar: varid: 17145 | x: 17140 | x[i]: 17145
Newvar: varid: 17146 | x: 17140 | x[i]: 17146
Newvar: varid: 17147 | x: 17140 | x[i]: 17147
Newvar: varid: 17148 | x: 17140 | x[i]: 17148
Newvar: varid: 17149 | x: 17140 | x[i]: 17149
Newvar: varid: 17150 | x: 17140 | x[i]: 17150
Newvar: varid: 17151 | x: 17140 | x[i]: 17151
Newvar: varid: 17152 | x: 17140 | x[i]: 17152
Newvar: varid: 17153 | x: 17140 | x[i]: 17153
Newvar: varid: 17154 | x: 17140 | x[i]: 17154
Newvar: varid: 17155 | x: 17140 | x[i]: 17155
Newvar: varid: 17156 | x: 17140 | x[i]: 17156
Newvar: varid: 17157 | x: 17140 | x[i]: 17157
Newvar: varid: 17158 | x: 17140 | x[i]: 17158
Newvar: varid: 17159 | x: 17140 | x[i]: 17159
Newvar: varid: 17160 | x: 17140 | x[i]: 17160
Newvar: varid: 17161 | x: 17140 | x[i]: 17161
Newvar: varid: 17162 | x: 17140 | x[i]: 17162
Newvar: varid: 17163 | x: 17140 | x[i]: 17163
Newvar: varid: 17164 | x: 17140 | x[i]: 17164
Newvar: varid: 17165 | x: 17140 | x[i]: 17165
Newvar: varid: 17166 | x: 17140 | x[i]: 17166
Newvar: varid: 17167 | x: 17140 | x[i]: 17167
Newvar: varid: 17168 | x: 17140 | x[i]: 17168
Newvar: varid: 17169 | x: 17140 | x[i]: 17169
Newvar: varid: 17170 | x: 17140 | x[i]: 17170
Defaults: z: 17108 | x: 1505 | y: 17044 | c: 17140
full: z: 17109 | x: 1506 | y: 17045 | c: 17140
xor3: z: 17139 | x: 1536 | y: 17075 | c: 17172
Newvar: varid: 17171 | x: 17171 | x[i]: 17171
Newvar: varid: 17172 | x: 17171 | x[i]: 17172
Newvar: varid: 17173 | x: 17171 | x[i]: 17173
Newvar: varid: 17174 | x: 17171 | x[i]: 17174
Newvar: varid: 17175 | x: 17171 | x[i]: 17175
Newvar: varid: 17176 | x: 17171 | x[i]: 17176
Newvar: varid: 17177 | x: 17171 | x[i]: 17177
Newvar: varid: 17178 | x: 17171 | x[i]: 17178
Newvar: varid: 17179 | x: 17171 | x[i]: 17179
Newvar: varid: 17180 | x: 17171 | x[i]: 17180
Newvar: varid: 17181 | x: 17171 | x[i]: 17181
Newvar: varid: 17182 | x: 17171 | x[i]: 17182
Newvar: varid: 17183 | x: 17171 | x[i]: 17183
Newvar: varid: 17184 | x: 17171 | x[i]: 17184
Newvar: varid: 17185 | x: 17171 | x[i]: 17185
Newvar: varid: 17186 | x: 17171 | x[i]: 17186
Newvar: varid: 17187 | x: 17171 | x[i]: 17187
Newvar: varid: 17188 | x: 17171 | x[i]: 17188
Newvar: varid: 17189 | x: 17171 | x[i]: 17189
Newvar: varid: 17190 | x: 17171 | x[i]: 17190
Newvar: varid: 17191 | x: 17171 | x[i]: 17191
Newvar: varid: 17192 | x: 17171 | x[i]: 17192
Newvar: varid: 17193 | x: 17171 | x[i]: 17193
Newvar: varid: 17194 | x: 17171 | x[i]: 17194
Newvar: varid: 17195 | x: 17171 | x[i]: 17195
Newvar: varid: 17196 | x: 17171 | x[i]: 17196
Newvar: varid: 17197 | x: 17171 | x[i]: 17197
Newvar: varid: 17198 | x: 17171 | x[i]: 17198
Newvar: varid: 17199 | x: 17171 | x[i]: 17199
Newvar: varid: 17200 | x: 17171 | x[i]: 17200
Newvar: varid: 17201 | x: 17171 | x[i]: 17201
Newvar: varid: 17202 | x: 17171 | x[i]: 17202
Newvar: varid: 17203 | x: 17203 | x[i]: 17203
Newvar: varid: 17204 | x: 17203 | x[i]: 17204
Newvar: varid: 17205 | x: 17203 | x[i]: 17205
Newvar: varid: 17206 | x: 17203 | x[i]: 17206
Newvar: varid: 17207 | x: 17203 | x[i]: 17207
Newvar: varid: 17208 | x: 17203 | x[i]: 17208
Newvar: varid: 17209 | x: 17203 | x[i]: 17209
Newvar: varid: 17210 | x: 17203 | x[i]: 17210
Newvar: varid: 17211 | x: 17203 | x[i]: 17211
Newvar: varid: 17212 | x: 17203 | x[i]: 17212
Newvar: varid: 17213 | x: 17203 | x[i]: 17213
Newvar: varid: 17214 | x: 17203 | x[i]: 17214
Newvar: varid: 17215 | x: 17203 | x[i]: 17215
Newvar: varid: 17216 | x: 17203 | x[i]: 17216
Newvar: varid: 17217 | x: 17203 | x[i]: 17217
Newvar: varid: 17218 | x: 17203 | x[i]: 17218
Newvar: varid: 17219 | x: 17203 | x[i]: 17219
Newvar: varid: 17220 | x: 17203 | x[i]: 17220
Newvar: varid: 17221 | x: 17203 | x[i]: 17221
Newvar: varid: 17222 | x: 17203 | x[i]: 17222
Newvar: varid: 17223 | x: 17203 | x[i]: 17223
Newvar: varid: 17224 | x: 17203 | x[i]: 17224
Newvar: varid: 17225 | x: 17203 | x[i]: 17225
Newvar: varid: 17226 | x: 17203 | x[i]: 17226
Newvar: varid: 17227 | x: 17203 | x[i]: 17227
Newvar: varid: 17228 | x: 17203 | x[i]: 17228
Newvar: varid: 17229 | x: 17203 | x[i]: 17229
Newvar: varid: 17230 | x: 17203 | x[i]: 17230
Newvar: varid: 17231 | x: 17203 | x[i]: 17231
Newvar: varid: 17232 | x: 17203 | x[i]: 17232
Newvar: varid: 17233 | x: 17203 | x[i]: 17233
Defaults: z: 17171 | x: 1793 | y: 17076 | c: 17203
full: z: 17172 | x: 1794 | y: 17077 | c: 17203
xor3: z: 17202 | x: 1824 | y: 17107 | c: 17235
Newvar: varid: 17234 | x: 17234 | x[i]: 17234
Newvar: varid: 17235 | x: 17234 | x[i]: 17235
Newvar: varid: 17236 | x: 17234 | x[i]: 17236
Newvar: varid: 17237 | x: 17234 | x[i]: 17237
Newvar: varid: 17238 | x: 17234 | x[i]: 17238
Newvar: varid: 17239 | x: 17234 | x[i]: 17239
Newvar: varid: 17240 | x: 17234 | x[i]: 17240
Newvar: varid: 17241 | x: 17234 | x[i]: 17241
Newvar: varid: 17242 | x: 17234 | x[i]: 17242
Newvar: varid: 17243 | x: 17234 | x[i]: 17243
Newvar: varid: 17244 | x: 17234 | x[i]: 17244
Newvar: varid: 17245 | x: 17234 | x[i]: 17245
Newvar: varid: 17246 | x: 17234 | x[i]: 17246
Newvar: varid: 17247 | x: 17234 | x[i]: 17247
Newvar: varid: 17248 | x: 17234 | x[i]: 17248
Newvar: varid: 17249 | x: 17234 | x[i]: 17249
Newvar: varid: 17250 | x: 17234 | x[i]: 17250
Newvar: varid: 17251 | x: 17234 | x[i]: 17251
Newvar: varid: 17252 | x: 17234 | x[i]: 17252
Newvar: varid: 17253 | x: 17234 | x[i]: 17253
Newvar: varid: 17254 | x: 17234 | x[i]: 17254
Newvar: varid: 17255 | x: 17234 | x[i]: 17255
Newvar: varid: 17256 | x: 17234 | x[i]: 17256
Newvar: varid: 17257 | x: 17234 | x[i]: 17257
Newvar: varid: 17258 | x: 17234 | x[i]: 17258
Newvar: varid: 17259 | x: 17234 | x[i]: 17259
Newvar: varid: 17260 | x: 17234 | x[i]: 17260
Newvar: varid: 17261 | x: 17234 | x[i]: 17261
Newvar: varid: 17262 | x: 17234 | x[i]: 17262
Newvar: varid: 17263 | x: 17234 | x[i]: 17263
Newvar: varid: 17264 | x: 17234 | x[i]: 17264
Defaults: z: 2017 | x: 17108 | y: 17171 | c: 17234
full: z: 2018 | x: 17109 | y: 17172 | c: 17234
xor3: z: 2048 | x: 17139 | y: 17202 | c: 17266
K assign...
Newvar: varid: 17265 | x: 17265 | x[i]: 17265
Newvar: varid: 17266 | x: 17265 | x[i]: 17266
Newvar: varid: 17267 | x: 17265 | x[i]: 17267
Newvar: varid: 17268 | x: 17265 | x[i]: 17268
Newvar: varid: 17269 | x: 17265 | x[i]: 17269
Newvar: varid: 17270 | x: 17265 | x[i]: 17270
Newvar: varid: 17271 | x: 17265 | x[i]: 17271
Newvar: varid: 17272 | x: 17265 | x[i]: 17272
Newvar: varid: 17273 | x: 17265 | x[i]: 17273
Newvar: varid: 17274 | x: 17265 | x[i]: 17274
Newvar: varid: 17275 | x: 17265 | x[i]: 17275
Newvar: varid: 17276 | x: 17265 | x[i]: 17276
Newvar: varid: 17277 | x: 17265 | x[i]: 17277
Newvar: varid: 17278 | x: 17265 | x[i]: 17278
Newvar: varid: 17279 | x: 17265 | x[i]: 17279
Newvar: varid: 17280 | x: 17265 | x[i]: 17280
Newvar: varid: 17281 | x: 17265 | x[i]: 17281
Newvar: varid: 17282 | x: 17265 | x[i]: 17282
Newvar: varid: 17283 | x: 17265 | x[i]: 17283
Newvar: varid: 17284 | x: 17265 | x[i]: 17284
Newvar: varid: 17285 | x: 17265 | x[i]: 17285
Newvar: varid: 17286 | x: 17265 | x[i]: 17286
Newvar: varid: 17287 | x: 17265 | x[i]: 17287
Newvar: varid: 17288 | x: 17265 | x[i]: 17288
Newvar: varid: 17289 | x: 17265 | x[i]: 17289
Newvar: varid: 17290 | x: 17265 | x[i]: 17290
Newvar: varid: 17291 | x: 17265 | x[i]: 17291
Newvar: varid: 17292 | x: 17265 | x[i]: 17292
Newvar: varid: 17293 | x: 17265 | x[i]: 17293
Newvar: varid: 17294 | x: 17265 | x[i]: 17294
Newvar: varid: 17295 | x: 17265 | x[i]: 17295
Newvar: varid: 17296 | x: 17265 | x[i]: 17296
fixed: z[i]: 17265 | z: 17265 | value: 1116352408 | Varid: -17265
fixed: z[i]: 17266 | z: 17265 | value: 1116352408 | Varid: -17266
fixed: z[i]: 17267 | z: 17265 | value: 1116352408 | Varid: -17267
fixed: z[i]: 17268 | z: 17265 | value: 1116352408 | Varid: 17268
fixed: z[i]: 17269 | z: 17265 | value: 1116352408 | Varid: 17269
fixed: z[i]: 17270 | z: 17265 | value: 1116352408 | Varid: -17270
fixed: z[i]: 17271 | z: 17265 | value: 1116352408 | Varid: -17271
fixed: z[i]: 17272 | z: 17265 | value: 1116352408 | Varid: 17272
fixed: z[i]: 17273 | z: 17265 | value: 1116352408 | Varid: 17273
fixed: z[i]: 17274 | z: 17265 | value: 1116352408 | Varid: 17274
fixed: z[i]: 17275 | z: 17265 | value: 1116352408 | Varid: 17275
fixed: z[i]: 17276 | z: 17265 | value: 1116352408 | Varid: 17276
fixed: z[i]: 17277 | z: 17265 | value: 1116352408 | Varid: -17277
fixed: z[i]: 17278 | z: 17265 | value: 1116352408 | Varid: 17278
fixed: z[i]: 17279 | z: 17265 | value: 1116352408 | Varid: -17279
fixed: z[i]: 17280 | z: 17265 | value: 1116352408 | Varid: -17280
fixed: z[i]: 17281 | z: 17265 | value: 1116352408 | Varid: -17281
fixed: z[i]: 17282 | z: 17265 | value: 1116352408 | Varid: 17282
fixed: z[i]: 17283 | z: 17265 | value: 1116352408 | Varid: -17283
fixed: z[i]: 17284 | z: 17265 | value: 1116352408 | Varid: 17284
fixed: z[i]: 17285 | z: 17265 | value: 1116352408 | Varid: -17285
fixed: z[i]: 17286 | z: 17265 | value: 1116352408 | Varid: -17286
fixed: z[i]: 17287 | z: 17265 | value: 1116352408 | Varid: -17287
fixed: z[i]: 17288 | z: 17265 | value: 1116352408 | Varid: 17288
fixed: z[i]: 17289 | z: 17265 | value: 1116352408 | Varid: -17289
fixed: z[i]: 17290 | z: 17265 | value: 1116352408 | Varid: 17290
fixed: z[i]: 17291 | z: 17265 | value: 1116352408 | Varid: -17291
fixed: z[i]: 17292 | z: 17265 | value: 1116352408 | Varid: -17292
fixed: z[i]: 17293 | z: 17265 | value: 1116352408 | Varid: -17293
fixed: z[i]: 17294 | z: 17265 | value: 1116352408 | Varid: -17294
fixed: z[i]: 17295 | z: 17265 | value: 1116352408 | Varid: 17295
fixed: z[i]: 17296 | z: 17265 | value: 1116352408 | Varid: -17296
Newvar: varid: 17297 | x: 17297 | x[i]: 17297
Newvar: varid: 17298 | x: 17297 | x[i]: 17298
Newvar: varid: 17299 | x: 17297 | x[i]: 17299
Newvar: varid: 17300 | x: 17297 | x[i]: 17300
Newvar: varid: 17301 | x: 17297 | x[i]: 17301
Newvar: varid: 17302 | x: 17297 | x[i]: 17302
Newvar: varid: 17303 | x: 17297 | x[i]: 17303
Newvar: varid: 17304 | x: 17297 | x[i]: 17304
Newvar: varid: 17305 | x: 17297 | x[i]: 17305
Newvar: varid: 17306 | x: 17297 | x[i]: 17306
Newvar: varid: 17307 | x: 17297 | x[i]: 17307
Newvar: varid: 17308 | x: 17297 | x[i]: 17308
Newvar: varid: 17309 | x: 17297 | x[i]: 17309
Newvar: varid: 17310 | x: 17297 | x[i]: 17310
Newvar: varid: 17311 | x: 17297 | x[i]: 17311
Newvar: varid: 17312 | x: 17297 | x[i]: 17312
Newvar: varid: 17313 | x: 17297 | x[i]: 17313
Newvar: varid: 17314 | x: 17297 | x[i]: 17314
Newvar: varid: 17315 | x: 17297 | x[i]: 17315
Newvar: varid: 17316 | x: 17297 | x[i]: 17316
Newvar: varid: 17317 | x: 17297 | x[i]: 17317
Newvar: varid: 17318 | x: 17297 | x[i]: 17318
Newvar: varid: 17319 | x: 17297 | x[i]: 17319
Newvar: varid: 17320 | x: 17297 | x[i]: 17320
Newvar: varid: 17321 | x: 17297 | x[i]: 17321
Newvar: varid: 17322 | x: 17297 | x[i]: 17322
Newvar: varid: 17323 | x: 17297 | x[i]: 17323
Newvar: varid: 17324 | x: 17297 | x[i]: 17324
Newvar: varid: 17325 | x: 17297 | x[i]: 17325
Newvar: varid: 17326 | x: 17297 | x[i]: 17326
Newvar: varid: 17327 | x: 17297 | x[i]: 17327
Newvar: varid: 17328 | x: 17297 | x[i]: 17328
fixed: z[i]: 17297 | z: 17297 | value: 1899447441 | Varid: 17297
fixed: z[i]: 17298 | z: 17297 | value: 1899447441 | Varid: -17298
fixed: z[i]: 17299 | z: 17297 | value: 1899447441 | Varid: -17299
fixed: z[i]: 17300 | z: 17297 | value: 1899447441 | Varid: -17300
fixed: z[i]: 17301 | z: 17297 | value: 1899447441 | Varid: 17301
fixed: z[i]: 17302 | z: 17297 | value: 1899447441 | Varid: -17302
fixed: z[i]: 17303 | z: 17297 | value: 1899447441 | Varid: -17303
fixed: z[i]: 17304 | z: 17297 | value: 1899447441 | Varid: 17304
fixed: z[i]: 17305 | z: 17297 | value: 1899447441 | Varid: -17305
fixed: z[i]: 17306 | z: 17297 | value: 1899447441 | Varid: -17306
fixed: z[i]: 17307 | z: 17297 | value: 1899447441 | Varid: 17307
fixed: z[i]: 17308 | z: 17297 | value: 1899447441 | Varid: -17308
fixed: z[i]: 17309 | z: 17297 | value: 1899447441 | Varid: -17309
fixed: z[i]: 17310 | z: 17297 | value: 1899447441 | Varid: -17310
fixed: z[i]: 17311 | z: 17297 | value: 1899447441 | Varid: 17311
fixed: z[i]: 17312 | z: 17297 | value: 1899447441 | Varid: -17312
fixed: z[i]: 17313 | z: 17297 | value: 1899447441 | Varid: 17313
fixed: z[i]: 17314 | z: 17297 | value: 1899447441 | Varid: 17314
fixed: z[i]: 17315 | z: 17297 | value: 1899447441 | Varid: 17315
fixed: z[i]: 17316 | z: 17297 | value: 1899447441 | Varid: -17316
fixed: z[i]: 17317 | z: 17297 | value: 1899447441 | Varid: 17317
fixed: z[i]: 17318 | z: 17297 | value: 1899447441 | Varid: 17318
fixed: z[i]: 17319 | z: 17297 | value: 1899447441 | Varid: -17319
fixed: z[i]: 17320 | z: 17297 | value: 1899447441 | Varid: -17320
fixed: z[i]: 17321 | z: 17297 | value: 1899447441 | Varid: 17321
fixed: z[i]: 17322 | z: 17297 | value: 1899447441 | Varid: -17322
fixed: z[i]: 17323 | z: 17297 | value: 1899447441 | Varid: -17323
fixed: z[i]: 17324 | z: 17297 | value: 1899447441 | Varid: -17324
fixed: z[i]: 17325 | z: 17297 | value: 1899447441 | Varid: 17325
fixed: z[i]: 17326 | z: 17297 | value: 1899447441 | Varid: 17326
fixed: z[i]: 17327 | z: 17297 | value: 1899447441 | Varid: 17327
fixed: z[i]: 17328 | z: 17297 | value: 1899447441 | Varid: -17328
Newvar: varid: 17329 | x: 17329 | x[i]: 17329
Newvar: varid: 17330 | x: 17329 | x[i]: 17330
Newvar: varid: 17331 | x: 17329 | x[i]: 17331
Newvar: varid: 17332 | x: 17329 | x[i]: 17332
Newvar: varid: 17333 | x: 17329 | x[i]: 17333
Newvar: varid: 17334 | x: 17329 | x[i]: 17334
Newvar: varid: 17335 | x: 17329 | x[i]: 17335
Newvar: varid: 17336 | x: 17329 | x[i]: 17336
Newvar: varid: 17337 | x: 17329 | x[i]: 17337
Newvar: varid: 17338 | x: 17329 | x[i]: 17338
Newvar: varid: 17339 | x: 17329 | x[i]: 17339
Newvar: varid: 17340 | x: 17329 | x[i]: 17340
Newvar: varid: 17341 | x: 17329 | x[i]: 17341
Newvar: varid: 17342 | x: 17329 | x[i]: 17342
Newvar: varid: 17343 | x: 17329 | x[i]: 17343
Newvar: varid: 17344 | x: 17329 | x[i]: 17344
Newvar: varid: 17345 | x: 17329 | x[i]: 17345
Newvar: varid: 17346 | x: 17329 | x[i]: 17346
Newvar: varid: 17347 | x: 17329 | x[i]: 17347
Newvar: varid: 17348 | x: 17329 | x[i]: 17348
Newvar: varid: 17349 | x: 17329 | x[i]: 17349
Newvar: varid: 17350 | x: 17329 | x[i]: 17350
Newvar: varid: 17351 | x: 17329 | x[i]: 17351
Newvar: varid: 17352 | x: 17329 | x[i]: 17352
Newvar: varid: 17353 | x: 17329 | x[i]: 17353
Newvar: varid: 17354 | x: 17329 | x[i]: 17354
Newvar: varid: 17355 | x: 17329 | x[i]: 17355
Newvar: varid: 17356 | x: 17329 | x[i]: 17356
Newvar: varid: 17357 | x: 17329 | x[i]: 17357
Newvar: varid: 17358 | x: 17329 | x[i]: 17358
Newvar: varid: 17359 | x: 17329 | x[i]: 17359
Newvar: varid: 17360 | x: 17329 | x[i]: 17360
fixed: z[i]: 17329 | z: 17329 | value: -1245643825 | Varid: 17329
fixed: z[i]: 17330 | z: 17329 | value: -1245643825 | Varid: 17330
fixed: z[i]: 17331 | z: 17329 | value: -1245643825 | Varid: 17331
fixed: z[i]: 17332 | z: 17329 | value: -1245643825 | Varid: 17332
fixed: z[i]: 17333 | z: 17329 | value: -1245643825 | Varid: -17333
fixed: z[i]: 17334 | z: 17329 | value: -1245643825 | Varid: -17334
fixed: z[i]: 17335 | z: 17329 | value: -1245643825 | Varid: 17335
fixed: z[i]: 17336 | z: 17329 | value: -1245643825 | Varid: 17336
fixed: z[i]: 17337 | z: 17329 | value: -1245643825 | Varid: 17337
fixed: z[i]: 17338 | z: 17329 | value: -1245643825 | Varid: 17338
fixed: z[i]: 17339 | z: 17329 | value: -1245643825 | Varid: -17339
fixed: z[i]: 17340 | z: 17329 | value: -1245643825 | Varid: 17340
fixed: z[i]: 17341 | z: 17329 | value: -1245643825 | Varid: 17341
fixed: z[i]: 17342 | z: 17329 | value: -1245643825 | Varid: 17342
fixed: z[i]: 17343 | z: 17329 | value: -1245643825 | Varid: 17343
fixed: z[i]: 17344 | z: 17329 | value: -1245643825 | Varid: 17344
fixed: z[i]: 17345 | z: 17329 | value: -1245643825 | Varid: -17345
fixed: z[i]: 17346 | z: 17329 | value: -1245643825 | Varid: -17346
fixed: z[i]: 17347 | z: 17329 | value: -1245643825 | Varid: -17347
fixed: z[i]: 17348 | z: 17329 | value: -1245643825 | Varid: -17348
fixed: z[i]: 17349 | z: 17329 | value: -1245643825 | Varid: -17349
fixed: z[i]: 17350 | z: 17329 | value: -1245643825 | Varid: -17350
fixed: z[i]: 17351 | z: 17329 | value: -1245643825 | Varid: 17351
fixed: z[i]: 17352 | z: 17329 | value: -1245643825 | Varid: 17352
fixed: z[i]: 17353 | z: 17329 | value: -1245643825 | Varid: 17353
fixed: z[i]: 17354 | z: 17329 | value: -1245643825 | Varid: -17354
fixed: z[i]: 17355 | z: 17329 | value: -1245643825 | Varid: 17355
fixed: z[i]: 17356 | z: 17329 | value: -1245643825 | Varid: -17356
fixed: z[i]: 17357 | z: 17329 | value: -1245643825 | Varid: 17357
fixed: z[i]: 17358 | z: 17329 | value: -1245643825 | Varid: 17358
fixed: z[i]: 17359 | z: 17329 | value: -1245643825 | Varid: -17359
fixed: z[i]: 17360 | z: 17329 | value: -1245643825 | Varid: 17360
Newvar: varid: 17361 | x: 17361 | x[i]: 17361
Newvar: varid: 17362 | x: 17361 | x[i]: 17362
Newvar: varid: 17363 | x: 17361 | x[i]: 17363
Newvar: varid: 17364 | x: 17361 | x[i]: 17364
Newvar: varid: 17365 | x: 17361 | x[i]: 17365
Newvar: varid: 17366 | x: 17361 | x[i]: 17366
Newvar: varid: 17367 | x: 17361 | x[i]: 17367
Newvar: varid: 17368 | x: 17361 | x[i]: 17368
Newvar: varid: 17369 | x: 17361 | x[i]: 17369
Newvar: varid: 17370 | x: 17361 | x[i]: 17370
Newvar: varid: 17371 | x: 17361 | x[i]: 17371
Newvar: varid: 17372 | x: 17361 | x[i]: 17372
Newvar: varid: 17373 | x: 17361 | x[i]: 17373
Newvar: varid: 17374 | x: 17361 | x[i]: 17374
Newvar: varid: 17375 | x: 17361 | x[i]: 17375
Newvar: varid: 17376 | x: 17361 | x[i]: 17376
Newvar: varid: 17377 | x: 17361 | x[i]: 17377
Newvar: varid: 17378 | x: 17361 | x[i]: 17378
Newvar: varid: 17379 | x: 17361 | x[i]: 17379
Newvar: varid: 17380 | x: 17361 | x[i]: 17380
Newvar: varid: 17381 | x: 17361 | x[i]: 17381
*/
