package graphs

import (
	"log"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
)

// valid returns true if m is a valid mapping, from sub node name to graph node
// name, for an isomorphism of sub in graph considering all nodes and edges
// except predecessors of entry and successors of exit.
func valid(graph *dot.Graph, sub *SubGraph, m map[string]string) bool {
	if len(m) != len(sub.Nodes.Nodes) {
		return false
	}
	for sname, gname := range m {
		s, ok := sub.Nodes.Lookup[sname]
		if !ok {
			err := errutil.Newf("unable to locate node %q in sub", sname)
			log.Println(err)
			return false
		}
		g, ok := graph.Nodes.Lookup[gname]
		if !ok {
			err := errutil.Newf("unable to locate node %q in graph", gname)
			log.Println(err)
			return false
		}

		// Verify predecessors.
		if s.Name != sub.entry {
			for _, spred := range s.Preds {
				found := false
				for _, gpred := range g.Preds {
					if gpred.Name == m[spred.Name] {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}

		// Verify successors.
		if s.Name != sub.exit {
			for _, ssucc := range s.Succs {
				found := false
				for _, gsucc := range g.Succs {
					if gsucc.Name == m[ssucc.Name] {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}
	}

	// Isomorphism found!
	return true
}
