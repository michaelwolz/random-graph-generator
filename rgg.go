// Author: Michael Wolz

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type graph struct {
	vertices  []int
	adjMatrix [][]uint8
}

func (g *graph) init(v int) {
	//add vertices
	for i := 0; i < v; i++ {
		g.addVertex(i)
	}

	//initialize adjacency matrix
	g.adjMatrix = make([][]uint8, v-1)
	for i := range g.adjMatrix {
		g.adjMatrix[i] = make([]uint8, v-1)
	}
}

func (g *graph) addVertex(ID int) {
	g.vertices = append(g.vertices, ID)
}

func (g *graph) addEdge(v1, v2 int) {
	v1, v2 = minMax(v1, v2)
	g.adjMatrix[v1][v2] = 1
}

func (g *graph) addRandomEdge() {
	// worst part :/
	var v = len(g.vertices)
	var v1, v2 int

	for v1 == v2 {
		v1 = rand.Intn(v - 1)
		v2 = rand.Intn(v - 1)
	}

	v1, v2 = minMax(v1, v2)
	fmt.Println(v1, v2)
	if g.adjMatrix[v1][v2] == 0 {
		g.adjMatrix[v1][v2] = 1
	} else {
		g.addRandomEdge()
	}
}

func (g *graph) printAdjMatrix() {
	fmt.Print("\n##### GRAPH ADJACENCY MATRIX #####\n\n")
	for i := 0; i < len(g.vertices)-1; i++ {
		fmt.Printf("(%d) %v\n", i+1, g.adjMatrix[i])
	}
	fmt.Print("\n##################################\n\n")
}

func minMax(v1, v2 int) (int, int) {
	//v1 - 1, because we don't have a first row! (lower triangular matrix)
	if v1 > v2 {
		return v1 - 1, v2
	}
	return v2 - 1, v1
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

//needed to define this function, because json.Marschal messes up uint8 values
func (g *graph) MarshalJSON() ([]byte, error) {
	var array string
	if g.adjMatrix == nil {
		array = "null"
	} else {
		array = strings.Join(strings.Fields(fmt.Sprintf("%d", g.adjMatrix)), ",")
	}
	jsonResult := fmt.Sprintf(`{"adjMatrix":%s}`, array)
	return []byte(jsonResult), nil
}

var maxEdges int

// ######################

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

	//seed the pseudo-rand generator
	rand.Seed(time.Now().UnixNano())

	//build graph
	g := buildRandomGraph(v, e)

	//ouput adjList
	g.printAdjMatrix()

	//write graph to JSON-file
	g.generateJSONGraph()
}

func buildRandomGraph(v, e int) graph {
	var g = graph{}
	g.init(v)

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
		g.addRandomEdge()
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
