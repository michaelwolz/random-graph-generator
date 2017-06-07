// Author: Michael Wolz

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type graph struct {
	Vertices []vertex `json:"vertices"`
}

type vertex struct {
	ID    int
	Edges []int `json:"edges"`
}

func (g *graph) addVertex(ID int) {
	g.Vertices = append(g.Vertices, vertex{ID, nil})
}

func (g *graph) addEdge(v1, v2 int) {
	g.Vertices[v1].Edges = append(g.Vertices[v1].Edges, v2)
	g.Vertices[v2].Edges = append(g.Vertices[v2].Edges, v1)
}

func (g *graph) traverseGraph() {
	fmt.Print("\n##### GRAPH #####\n\n")
	for i, v := range g.Vertices {
		fmt.Printf("(%d) -> %v\n", i, v.Edges)
	}
	fmt.Print("\n##################\n\n")
}

func (g *graph) generateJSONGraph() {
	json, _ := json.Marshal(g)

	f, err := os.Create("graph.json")
	check(err)

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(json))
	check(err)
	w.Flush()

	fmt.Print("JSON-Data written to file: ./graph.json\n\n")
}

var maxEdges int

func main() {
	//init
	var args = os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: rgg vertices edges")
		os.Exit(1)
	}

	v := argParse(args[0])
	e := argParse(args[1])
	maxEdges = v * (v - 1) / 2

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
	g.generateJSONGraph()
}

func buildRandomGraph(v, e int) graph {
	var g = graph{}
	for i := 0; i < v; i++ {
		g.addVertex(i)
	}

	if e == maxEdges {
		//if amount of edges equals maxEdges, just connect ALL vertices.
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

	//randomly add the remaining edges to the graph
	remaining := e - v + 1
	for remaining > 0 {
		remaining--
	}
}

func argParse(arg string) int {
	res, err := strconv.Atoi(arg)
	check(err)

	return res
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
