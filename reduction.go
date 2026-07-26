package main

/* reduction.go


*/

type Graph struct {
	NodeCount int
	Labels    map[string]int
	Nodes     []Node
	Edges     []Edge
}

type Node struct {
	Label    string
	Variable string
	ID       int
}

type Edge struct {
	S Node
	D Node
}

/* transform cnf to graph problem */
func transform(fx Formula) Graph {
	g := Graph{}
	g.AddNode("Ground")
	g.AddNode("False")
	g.Connect(g.Nodes[0], g.Nodes[1])
	for i := 0; i < fx.VarCount(); i++ {
		g.AddVarGadget(i + 1)
	}
	for i := 0; i < len(fx.Clauses); i++ {
		g.AddClauseGadget(fx.Clauses[i])
	}
	return g
}

func (g *Graph) AddClauseGadget(c Clause) {
	// g.AddVarGadget()
	/*
		1.
		2.
		3.
	*/
	j := len(c.Literals) * 2
	for i := 0; i < j; i++ {
		g.AddNode(itoa(0))
		switch i {
		case 0:

		default:
			g.Connect(g.Nodes[g.NodeCount], g.Nodes[g.NodeCount-1])
		}
		if i == j {

		}
	}
	g.Connect(g.Nodes[g.NodeCount], g.Nodes[0])
	g.Connect(g.Nodes[g.NodeCount], g.Nodes[1])
}

/* connect node n to nodes n2... */
func (g *Graph) Connect(src Node, dsts ...Node) {
	for i := 0; i < len(dsts); i++ {
		g.AddEdge(src, dsts[i])
	}
}

/* Create a variable gadget */
func (g *Graph) AddVarGadget(v int) {
	g.AddNode(itoa(v))
	g.Connect(g.Nodes[g.NodeCount], g.Nodes[0])
	g.AddNode(itoa(-v))
	g.Connect(g.Nodes[g.NodeCount], g.Nodes[0])
	g.Connect(g.Nodes[g.NodeCount], g.Nodes[g.NodeCount-1])
}

/* add node to graph */
func (g *Graph) AddNode(l string) {
	n := Node{}
	n.ID = g.NodeCount
	g.NodeCount += 1
	n.Label = l
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph) AddEdge(n Node, n2 Node) {
	e := Edge{}
	e.S = n
	e.D = n2
	g.Edges = append(g.Edges, e)
}

/* append graph gx to g */
func (g *Graph) AppendGraph(gx Graph) {
	idx := g.NodeCount
	_ = idx
	for i := 0; i < gx.NodeCount; i++ {

	}
}

func (g *Graph)Walk() {

}





