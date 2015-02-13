// iso is a tool which locates subgraph isomorphisms in graphs.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
Usage: iso GRAPH.dot SUB.dot
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
	graphPath, subPath := flag.Arg(0), flag.Arg(1)
	err := iso(graphPath, subPath)
	if err != nil {
		log.Fatalln(err)
	}
}

// iso parses the provided graph and subgraph and tries to locate isomorphisms
// of the subgraph in the graph.
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
				fmt.Println("FOUND:", nodes[entry], "m:", m)
			}
		}
	} else {
		// Locate the first isomorphism of subgraph in graph.
		m, ok := graphs.Search(graph, sub)
		if ok {
			found = true
			fmt.Println("FOUND:", nodes[m[sub.Entry()]], "m:", m)
		}
	}
	if !found {
		fmt.Println("not found.")
	}

	return nil
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
