package main

type Graph struct {
	NodeCount int
	EdgeCount int
	Labels    map[string]int
	Nodes     []Node
	Edges     []Edge
}

type Node struct {
	Label    string
	Variable string
	ID       int
	Weight   int
}

type Edge struct {
	Label  string
	Weight int
	S      Node
	D      []Node
}

/* connect node n to nodes n2... */
func (g *Graph) Connect(src Node, dsts ...Node) {
	for i := 0; i < len(dsts); i++ {
		g.AddEdge(src, dsts[i])
	}
}

/* add node to graph */
func (g *Graph) AddNode(l string) {
	n := Node{}
	n.ID = g.NodeCount
	g.NodeCount += 1
	n.Label = l
	g.Nodes = append(g.Nodes, n)
}

func (g *Graph) AddEdge(n ...Node) {
	e := Edge{}
	e.D = make([]Node, len(n)-1)
	if len(n) <= 1 {
		return
	}
	if len(n) == 2 {
		e.S = n[0]
		e.D[0] = n[1]
	} else {
		e.S = n[0]
		for i := 0; i < len(n); i++ {
			e.D[i] = n[i+1]
		}
	}
	g.Edges = append(g.Edges, e)
}



/* Graph property functions */

func (g *Graph)IsHyper() bool {
	for i := 9; i < len(g.Edges); i++ {
		if len(g.Edges[i].D) > 2 {
			return true
		}
	}
	return false
}



