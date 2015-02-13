// iso is a tool which locates subgraph isomorphisms in graphs.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewrev/graphs"
)

// When flagAll is true, locate all isomorphisms of the subgraph in the graph.
var flagAll bool

func init() {
	flag.BoolVar(&flagAll, "all", true, "Locate all isomorphisms of SUB in GRAPH.")
	flag.Usage = usage
}

const use = `
Usage: iso [OPTION]... SUB.dot GRAPH.dot
Locates isomorphisms of the subgraph SUB in GRAPH.

Flags:`

func usage() {
	fmt.Fprintln(os.Stderr, use[1:])
	flag.PrintDefaults()
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}
	subPath, graphPath := flag.Arg(0), flag.Arg(1)
	err := iso(graphPath, subPath)
	if err != nil {
		log.Fatalln(err)
	}
}

// iso parses the provided graphs and tries to locate isomorphisms of the
// subgraph in the graph.
func iso(graphPath, subPath string) error {
	// Parse graphs.
	graph, err := parseGraph(graphPath)
	if err != nil {
		return errutil.Err(err)
	}
	sub, err := parseSubGraph(subPath)
	if err != nil {
		return errutil.Err(err)
	}
	nodes := graph.Nodes.Nodes

	found := false
	if flagAll {
		// Locate all isomorphisms of subgraph in graph.
		for entry := 0; entry < len(nodes); entry++ {
			m, ok := graphs.Isomorphism(graph, entry, sub)
			if ok {
				found = true
				printMapping(graph, sub, m)
			}
		}
	} else {
		// Locate the first isomorphism of subgraph in graph.
		m, ok := graphs.Search(graph, sub)
		if ok {
			found = true
			printMapping(graph, sub, m)
		}
	}
	if !found {
		fmt.Println("not found.")
	}

	return nil
}

// printMapping prints the mapping from sub node index to graph node index for
// an isomorphism of sub in graph.
func printMapping(graph *dot.Graph, sub *graphs.SubGraph, m map[int]int) {
	gnodes, snodes := graph.Nodes.Nodes, sub.Nodes.Nodes
	entry := m[sub.Entry()]
	var sidxs []int
	for sidx := range m {
		sidxs = append(sidxs, sidx)
	}
	sort.Ints(sidxs)
	fmt.Printf("Isomorphism found at node %q:\n", gnodes[entry].Name)
	for _, sidx := range sidxs {
		fmt.Printf("   %q=%q\n", snodes[sidx].Name, gnodes[m[sidx]].Name)
	}
}

// parseGraph parses the provided DOT file into a graph.
func parseGraph(path string) (*dot.Graph, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errutil.Err(err)
	}
	graph, err := dot.Read(buf)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return graph, nil
}

// parseSubGraph parses the provided DOT file into a subgraph with a dedicated
// entry and exit node. The entry and exit nodes are identified using the node
// "label" attribute, e.g.
//
//    digraph if {
//       A->B [label="true"]
//       A->C [label="false"]
//       B->C
//       A [label="entry"]
//       B
//       C [label="exit"]
//    }
func parseSubGraph(path string) (*graphs.SubGraph, error) {
	graph, err := parseGraph(path)
	if err != nil {
		return nil, errutil.Err(err)
	}
	return graphs.NewSubGraph(graph)
}
