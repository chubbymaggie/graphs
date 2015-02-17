// iso is a tool which locates subgraph isomorphisms in graphs.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewrev/graphs"
)

// When flagStart is a non-empty string, locate an isomorphism of the subgraph
// in the graph which starts at the given node.
var flagStart string

func init() {
	flag.StringVar(&flagStart, "start", "", "Locate an isomorphism of SUB in GRAPH which starts at the given node.")
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
	graph, err := dot.ParseFile(graphPath)
	if err != nil {
		return errutil.Err(err)
	}
	sub, err := graphs.ParseSubGraph(subPath)
	if err != nil {
		return errutil.Err(err)
	}

	// Locate isomorphisms.
	for i := 0; i < 1000; i++ {
		found := false
		if len(flagStart) > 0 {
			// Locate an isomorphism of sub in graph which starts at the node
			// specified by the "-start" flag.
			m, ok := graphs.Isomorphism(graph, flagStart, sub)
			if ok {
				found = true
				printMapping(graph, sub, m)
			}
		} else {
			// Locate all isomorphisms of sub in graph.
			var names []string
			for name := range graph.Nodes.Lookup {
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				m, ok := graphs.Isomorphism(graph, name, sub)
				if !ok {
					continue
				}
				found = true
				//printMapping(graph, sub, m)
				fmt.Println("found:", sorted(m))
			}
		}
		if !found {
			fmt.Println("not found.")
		}
	}

	return nil
}

func sorted(m map[string]string) string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var s string
	for _, key := range keys {
		s += fmt.Sprintf("%q:%q, ", key, m[key])
	}
	return s

}

// printMapping prints the mapping from sub node name to graph node name for an
// isomorphism of sub in graph.
func printMapping(graph *dot.Graph, sub *graphs.SubGraph, m map[string]string) {
	entry := m[sub.Entry()]
	var snames []string
	for sname := range m {
		snames = append(snames, sname)
	}
	sort.Strings(snames)
	fmt.Printf("Isomorphism found at node %q:\n", entry)
	for _, sname := range snames {
		fmt.Printf("   %q=%q\n", sname, m[sname])
	}
}
