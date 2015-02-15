// merge is a tool which merges subgraph isomorphisms in graphs into single
// nodes.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewrev/graphs"
)

var (
	// When flagAll is true, merge all isomorphisms of the subgraph in the graph
	// to single nodes.
	flagAll bool
	// When flagImage is true, generate an image representation of the CFG.
	flagImage bool
	// When flagQuiet is true, suppress non-error messages.
	flagQuiet bool
)

func init() {
	flag.BoolVar(&flagAll, "all", true, "Merge all isomorphisms of SUB in GRAPH to single nodes.")
	flag.BoolVar(&flagImage, "img", false, "Generate an image representation of the CFG.")
	flag.BoolVar(&flagQuiet, "q", false, "Suppress non-error messages.")
	flag.Usage = usage
}

const use = `
Usage: merge [OPTION]... SUB.dot GRAPH.dot
Merges isomorphisms of the subgraph SUB in GRAPH to single nodes.

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
	err := merge(graphPath, subPath)
	if err != nil {
		log.Fatalln(err)
	}
}

// merge parses the provided graphs and tries to merge isomorphisms of the
// subgraph in the graph into single nodes.
func merge(graphPath, subPath string) error {
	// Parse graphs.
	graph, err := parseGraph(graphPath)
	if err != nil {
		return errutil.Err(err)
	}
	sub, err := parseSubGraph(subPath)
	if err != nil {
		return errutil.Err(err)
	}

	// Merge isomorphisms of subgraph in graph into single nodes.
	found := false
	for entry := 0; entry < len(graph.Nodes.Nodes); entry++ {
		m, ok := graphs.Isomorphism(graph, entry, sub)
		if !ok {
			continue
		}
		found = true
		printMapping(graph, sub, m)

		entry, exit := graph.Nodes.Nodes[m[sub.Entry()]], graph.Nodes.Nodes[m[sub.Exit()]]
		err = graph.Merge(entry, exit, uniqName(sub.Name))
		if err != nil {
			return err
		}

		// Break after first merge.
		if !flagAll {
			break
		}
	}

	// TODO: Consider using os.Exit codes to signal that a subgraph was
	// successfully located and merged. This would enable loops in bash scripts.
	if found {
		err = dump(graph, "out")
		if err != nil {
			return errutil.Err(err)
		}
	} else {
		fmt.Println("SUB not present in GRAPH.")
	}

	return nil
}

// dump stores the graph as a DOT file and an image representation of the graph
// as a PNG file with filenames based on the given name.
func dump(graph *dot.Graph, name string) error {
	// Store graph to DOT file.
	dotPath := name + ".dot"
	if !flagQuiet {
		log.Printf("Creating: %q\n", dotPath)
	}
	err := ioutil.WriteFile(dotPath, []byte(graph.String()), 0644)
	if err != nil {
		return errutil.Err(err)
	}

	// Generate an image representation of the graph.
	if flagImage {
		pngPath := name + ".png"
		if !flagQuiet {
			log.Printf("Creating: %q\n", pngPath)
		}
		cmd := exec.Command("dot", "-Tpng", "-o", pngPath, dotPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return errutil.Err(err)
		}
	}

	return nil
}

// uniq maps from name to the next unused numeric suffix.
var uniq = make(map[string]int)

// uniqName returns name with a uniq numeric suffix.
func uniqName(name string) string {
	id := uniq[name]
	uniq[name]++
	return fmt.Sprintf("%s%d", name, id)
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
