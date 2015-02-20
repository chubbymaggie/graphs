package graphs

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
)

// Candidates locates a mapping from sub node name to graph node name candidates
// for an isomorphism of sub in graph which starts at the entry node.
func Candidates(graph *dot.Graph, entry string, sub *SubGraph) map[string]map[string]bool {
	g, ok := graph.Nodes.Lookup[entry]
	if !ok {
		err := errutil.Newf("unable to locate entry node %q in graph", entry)
		log.Println(err)
		return nil
	}
	s, ok := sub.Nodes.Lookup[sub.entry]
	if !ok {
		err := errutil.Newf("unable to locate entry node %q in sub", sub.entry)
		log.Println(err)
		return nil
	}
	if !isPotential(g, s, sub) {
		err := errutil.Newf("invalid entry node candidate %q; expected %d successors, got %d", g.Name, len(s.Succs), len(g.Succs))
		log.Println(err)
		return nil
	}
	c := map[string]map[string]bool{
		sub.entry: {
			entry: true,
		},
	}
	locate(s, g, sub, c)
	if len(c) != len(sub.Nodes.Nodes) {
		err := errutil.Newf("incomplete mapping; expected %d map entities, got %d", len(sub.Nodes.Nodes), len(c))
		log.Println(err)
		fmt.Println("### [ incomplete mapping ] ###")
		spew.Dump(c)
		fmt.Println("### [/ incomplete mapping ] ###")
		return nil
	}
	return c
}

// locate recursively locates potential node pairs (g and s) for an isomorphism
// of sub in graph and adds them to c, which is a mapping from sub node name to
// graph node name candidates.
func locate(g, s *dot.Node, sub *SubGraph, c map[string]map[string]bool) {
	// Early exit for impossible node pairs.
	if !isPotential(g, s, sub) {
		return
	}

	// Add g as a sub node name candidate.
	if s.Name != sub.entry {
		_, ok := c[s.Name]
		if !ok {
			c[s.Name] = make(map[string]bool)
		}
		c[s.Name][g.Name] = true
	}

	// TODO: Prevent inf cycles.
	for _, ssucc := range s.Succs {
		for _, gsucc := range g.Succs {
			locate(gsucc, ssucc, sub, c)
		}
	}
}

// isPotential returns true if the graph node g is a potential candidate for the
// sub node s, and false otherwise.
func isPotential(g, s *dot.Node, sub *SubGraph) bool {
	// Verify predecessors.
	if s.Name != sub.entry {
		if len(g.Preds) != len(s.Preds) {
			return false
		}
	}

	// Verify successors.
	if s.Name != sub.exit {
		if len(g.Succs) != len(s.Succs) {
			return false
		}
	}

	return true
}
