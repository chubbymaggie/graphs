// Package graphs implements subgraph isomorphism algorithms.
package graphs

import (
	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
)

// SubGraph represents a subgraph with a dedicated entry and exit node. Incoming
// edges to entry and outgoing edges from exit are ignored when searching for
// isomorphisms of the subgraph.
type SubGraph struct {
	*dot.Graph
	entry, exit int
}

// NewSubGraph returns a new subgraph based on graph with a dedicated entry and
// exit node. The entry and exit nodes are identified using the node "label"
// attribute, e.g.
//
//    digraph if {
//       A->B [label="true"]
//       A->C [label="false"]
//       B->C
//       A [label="entry"]
//       B
//       C [label="exit"]
//    }
func NewSubGraph(graph *dot.Graph) (*SubGraph, error) {
	sub := &SubGraph{Graph: graph}

	// Locate entry and exit nodes.
	var hasEntry, hasExit bool
	for _, node := range graph.Nodes.Nodes {
		label, ok := node.Attrs["label"]
		if !ok {
			continue
		}
		switch label {
		case "entry":
			if hasEntry {
				return nil, errutil.Newf("redefinition of entry node; previous index (%d), new index (%d)", sub.entry, node.Index)
			}
			sub.entry = node.Index
			hasEntry = true
		case "exit":
			if hasExit {
				return nil, errutil.Newf("redefinition of exit node; previous index (%d), new index (%d)", sub.exit, node.Index)
			}
			sub.exit = node.Index
			hasExit = true
		}
	}
	if !hasEntry {
		return nil, errutil.New("unable to locate entry node")
	}
	if !hasExit {
		return nil, errutil.New("unable to locate exit node")
	}

	return sub, nil
}

// Isomorph returns a mapping from sub node index to graph node index if there
// exists an isomorphism of sub in graph (starting at the entry node index). The
// boolean value is true if such a mapping could be located, and false
// otherwise.
func Isomorph(graph *dot.Graph, entry int, sub *SubGraph) (m map[int]int, ok bool) {
	panic("not yet implemented.")
}

// Search tries to locate an isomorphism of sub in graph. If successful it
// returns the mapping from sub node index to graph node index of the first
// isomorphism located. The boolean value is true if such a mapping could be
// located, and false otherwise.
func Search(graph *dot.Graph, sub *SubGraph) (m map[int]int, ok bool) {
	panic("not yet implemented.")
}
