package main

import "fmt"

/* reduction.go

*/

/* transform cnf to graph problem */
func transform(fx Formula) Graph {
	g := Graph{}
	g.AddNode("Ground")
	g.AddNode("False")
	g.Connect(g.Nodes[0], g.Nodes[1])
	for i := 0; i < fx.VarCount(); i++ {
		// g.AddVarGadget(i + 1)
	}
	// fmt.Println("Clauses:", fx.Clauses)
	for i := 0; i < len(fx.Clauses); i++ {
		g.AddClauseGadget(fx.Clauses[i])
	}
	return g
}

func (g *Graph) AddClauseGadget(c Clause) {
	if len(c.Literals) > 1 {
		fmt.Println("", len(c.Literals))
	}
	for i := 0; i < len(c.Literals); i++ {
		// v := c.Literals[i]
		// fmt.Println("Lit:", v)
	}
	// j := len(c.Literals) * 2
	// for i := 0; i < j; i++ {
	// 	g.AddNode(itoa(0))
	// 	switch i {
	// 	case 0:
	// 	default:
	// 		g.Connect(g.Nodes[g.Order], g.Nodes[g.Order-1])
	// 	}
	// 	if i == j {
	// 	}
	// }
	// g.Connect(g.Nodes[g.Order], g.Nodes[0])
	// g.Connect(g.Nodes[g.Order], g.Nodes[1])
}

/* Create a variable gadget */
func (g *Graph) AddVarGadget(v int) {
	g.AddNode(itoa(v))
	if len(g.Nodes) == g.Order {
		return
	}
	g.Connect(g.Nodes[g.Order], g.Nodes[0])
	g.AddNode(itoa(-v))
	g.Connect(g.Nodes[g.Order], g.Nodes[0])
	g.Connect(g.Nodes[g.Order], g.Nodes[g.Order-1])
}

/* append graph gx to g */
func (g *Graph) AppendGraph(gx Graph) {
	idx := g.Order
	_ = idx
	for i := 0; i < gx.Order; i++ {
		
	}
}

/*  */
func transformHyper(fx Formula) Graph {
	g := Graph{}
	//
	for i := 0; i < len(fx.Clauses); i++ {
		
	}
	for i := 0; i < fx.varCount; i++ {
		
	}
	return g
}




/* 

type Hgraph struct {
	X []clauses
	Y map[clause_id][]appearances_int
	Z [][]int //
	
	:
}

 */
