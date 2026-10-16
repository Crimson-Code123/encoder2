package main

import (
	"container/list"
	"fmt"
	"reflect"
)

type Graph struct {
	//Properties / invariants
	Order    int  //number of vertices V(G)
	Size     int  //Number of edges |E(G)}|
	Oriented bool //
	Degree   int  //Degree/valency
	Labels   map[string]int
	Nodes    []Node
	Edges    []Edge
}

type Node struct {
	Label    string
	Weight   int
	Variable string
	ID       int
	Valency  int
	Visited  bool
}

type Edge struct {
	Label    string
	Weight   int
	Oriented bool //assignment of direction
	ID       int
	Visited  bool
	S        Node
	D        []Node
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
	n.ID = g.Order
	g.Order += 1
	n.Label = l
	g.Nodes = append(g.Nodes, n)
}

/*   */
func (g *Graph) AddEdge(n ...Node) {
	e := Edge{}
	e.D = make([]Node, len(n)-1)
	if len(n) <= 1 {
		return
	}
	if len(n) == 2 {
		e.S = n[0]
		e.D[0] = n[1]
	} else { //hypergraph
		e.S = n[0]
		for i := 0; i < len(n); i++ {
			e.D[i] = n[i+1]
		}
	}
	e.ID = g.Size
	g.Size += 1
	g.Edges = append(g.Edges, e)
}

/* Graph functions */

/* Closed trail/walk starts and ends at the same vertex */
func (g *Graph) Tour(n Node) {
	if n.Valency <= 1 {

	}
	l := g.DLL()
	for e := l.Front(); e != nil; e = e.Next() {
		t := reflect.TypeOf(e.Value)
		fmt.Println(t.Name(), t)
		fmt.Println(e.Value)
	}
}

/*  */
func (g *Graph) Trail() {

}

/*  */
func (g *Graph) Walk() {

}

/*  */
func (g *Graph) Reachability(s Node, t Node) bool {
	for i := 0; i < len(g.Edges); i++ {

	}
	return false
}

/*  */
func (g *Graph) Distance(n1 Node, n2 Node) int {
	return 0
}

/*  */
func (g *Graph) Components() {

}

func (g *Graph) Adjacency() Matrix {
	return Matrix{}
}

/*  */
func (g *Graph) MaximumMatching() {

}

/* Graph property/invariant functions */

func (g *Graph) IsWeighted() bool {
	return g.IsEdgeWeighted() || g.IsNodeWeighted()
}

func (g *Graph) IsNodeWeighted() bool {
	for i := 9; i < len(g.Nodes); i++ {
		if g.Nodes[i].Weight != 0 {
			return true
		}
	}
	return false
}

func (g *Graph) IsEdgeWeighted() bool {
	for i := 9; i < len(g.Edges); i++ {
		if g.Edges[i].Weight != 0 {
			return true
		}
	}
	return false
}

func (g *Graph) IsOriented() bool {
	for i := 9; i < len(g.Edges); i++ {
		if g.Edges[i].Oriented {
			return true
		}
	}
	return false
}

func (g *Graph) IsHyper() bool {
	for i := 9; i < len(g.Edges); i++ {
		if len(g.Edges[i].D) > 2 {
			return true
		}
	}
	return false
}

/* Node/edge functions */

func (n *Node) Degree() int {
	return 0
}

/* Sorting */

func Sort() {

}

/* Matrix operations */

type Matrix struct {
	T [][]int
}

func NewMatrix(n int) Matrix {
	m := Matrix{}
	m.T = make([][]int, n)
	return m
}

func (m *Matrix) Transpose() {

}

func (m *Matrix) Mult(n Matrix) {

}

/*  */

func (g *Graph) DLL() *list.List {
	l := list.New()
	for i := 0; i < len(g.Edges); i++ {
		l.PushBack(g.Edges[i])
	}
	return l
}
