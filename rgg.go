// Author: Michael Wolz

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type graph struct {
	vertices []vertex
}

type vertex struct {
	edges []int
}

func (g *graph) addVertex() {
	g.vertices = append(g.vertices, vertex{})
}

func (g *graph) addEdge(v1, v2 int) {
	g.vertices[v1].edges = append(g.vertices[v1].edges, v2)
	g.vertices[v2].edges = append(g.vertices[v2].edges, v1)
}

func (g *graph) traverseGraph() {
	for i, v := range g.vertices {
		fmt.Printf("(%d) -> %v\n", i, v)
	}
}

func main() {
	//init
	var args = os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: rgg vertices edges")
		os.Exit(1)
	}

	v := argParse(args[0])
	e := argParse(args[1])
	maxEdges := v * (v - 1) / 2

	if e > maxEdges {
		fmt.Fprintf(os.Stderr, "error: max amount of edges in a graph with %d vertices is %d. Edge weight must be 1\n", v, maxEdges)
		os.Exit(1)
	}

	if e < v-1 {
		fmt.Fprintf(os.Stderr, "error: min amount of edges in a graph with %d vertices is %d. Edge weight must be 1\n", v, v-1)
		os.Exit(1)
	}

	//seed the pseudo-rand generators
	rand.Seed(time.Now().UnixNano())

	//build graph
	g := buildRandomGraph(v, e)

	//ouput adjList
	g.traverseGraph()
}

func buildRandomGraph(v, e int) graph {
	var g = graph{}
	for i := 0; i < v; i++ {
		g.addVertex()
	}

	if v-1 == e {
		//if amount of edges equals amount of vertices - 1, just connect ALL vertices. Much faster!
		for i := 0; i < v; i++ {
			for j := i + 1; j < v; j++ {
				g.addEdge(i, j)
			}
		}
	} else {
		distributeEdges(g, v, e)
	}

	return g
}

func distributeEdges(g graph, v, e int) {
	//connect all vertices with v-1 edges
	var vertexPermutation = rand.Perm(v)
	for i := 0; i < len(vertexPermutation)-1; i++ {
		g.addEdge(vertexPermutation[i], vertexPermutation[i+1])
	}

	//and add the remaining edges
	remaining := e - v + 1
	for remaining > 0 {
		remaining--
	}
}

func generateJSONOutput() {
}

func argParse(arg string) int {
	res, err := strconv.Atoi(arg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "vertices and edges must be numbers")
		panic(err)
	} else {
		return res
	}
}
