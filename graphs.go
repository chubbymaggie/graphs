// Package graphs implements subgraph isomorphism search algorithms.
package graphs

import (
	"fmt"

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
				return nil, errutil.Newf(`redefinition of "entry" node; previous name (%v), new name (%v)`, sub.Nodes.Nodes[sub.entry], graph.Nodes.Nodes[node.Index])
			}
			sub.entry = node.Index
			hasEntry = true
		case "exit":
			if hasExit {
				return nil, errutil.Newf(`redefinition of "exit" node; previous name (%d), new name (%d)`, sub.Nodes.Nodes[sub.exit], graph.Nodes.Nodes[node.Index])
			}
			sub.exit = node.Index
			hasExit = true
		}
	}
	if !hasEntry {
		return nil, errutil.New(`unable to locate node with "entry" label`)
	}
	if !hasExit {
		return nil, errutil.New(`unable to locate node with "exit" label`)
	}

	return sub, nil
}

// Entry returns the entry node index in the subgraph.
func (sub *SubGraph) Entry() int {
	return sub.entry
}

// Exit returns the exit node index in the subgraph.
func (sub *SubGraph) Exit() int {
	return sub.exit
}

// Search tries to locate an isomorphism of sub in graph. If successful it
// returns the mapping from sub node index to graph node index of the first
// isomorphism located. The boolean value is true if such a mapping could be
// located, and false otherwise.
func Search(graph *dot.Graph, sub *SubGraph) (m map[int]int, ok bool) {
	for entry := 0; entry < len(graph.Nodes.Nodes); entry++ {
		m, ok = Isomorphism(graph, entry, sub)
		if ok {
			return m, true
		}
	}
	return nil, false
}

// Isomorphism returns a mapping from sub node index to graph node index if
// there exists an isomorphism of sub in graph (starting at the entry node
// index). The boolean value is true if such a mapping could be located, and
// false otherwise.
func Isomorphism(graph *dot.Graph, entry int, sub *SubGraph) (m map[int]int, ok bool) {
	m = make(map[int]int)
	g := graph.Nodes.Nodes[entry]
	s := sub.Graph.Nodes.Nodes[sub.entry]
	if isIsomorphism(g, s, graph, sub, m) {
		return m, true
	}
	return nil, false
}

// isIsomorphism returns true if g is an isomorphism of s, where g is a node of
// graph, s is a node of sub and m is a mapping from sub node index to graph
// node index. Incoming edges to entry and outgoing edges from exit are ignored
// when searching for isomorphisms of sub.
func isIsomorphism(g, s *dot.Node, graph *dot.Graph, sub *SubGraph, m map[int]int) bool {
	// TODO: Check for loops?
	// TODO: Check for duplicate val in m and only add if not already present.
	fmt.Println("s.Index:", s.Index)
	fmt.Println("m:", m)

	// Create mapping from sub node index to graph node index by trying possible
	// candidates.
	if len(m) != len(sub.Nodes.Nodes) {
		if s.Index == sub.entry {
			// Add mapping for entry node.
			fmt.Println("ENTRY:", s.Index)
			if _, ok := m[s.Index]; ok {
				// TODO: How are graphs with entry points having outgoing edges to
				// themselves handled?
				panic("should only arrive here once")
			}
			m[s.Index] = g.Index
		} else {
			// Verify predecessor count.
			if len(s.Preds) != len(g.Preds) {
				return false
			}
		}

		// Add candidate mapping for successors.
		if s.Index != sub.exit {
			if len(s.Succs) != len(g.Succs) {
				return false
			}
			for _, ssucc := range s.Succs {
				fmt.Println("ssucc:", ssucc.Index)
				for _, gsucc := range g.Succs {
					fmt.Println("gsucc:", gsucc.Index)
					if _, ok := m[ssucc.Index]; !ok {
						m[ssucc.Index] = gsucc.Index
					}
					if isIsomorphism(gsucc, ssucc, graph, sub, m) {
						return true
					}
					delete(m, ssucc.Index)
				}
			}
		}
	}

	// Complete mapping?
	if len(m) != len(sub.Nodes.Nodes) {
		return false
	}

	// Check if m correctly maps the nodes of sub onto the nodes of graph.
	for sidx, gidx := range m {
		snode, gnode := sub.Nodes.Nodes[sidx], graph.Nodes.Nodes[gidx]

		// Check predecessors.
		if snode.Index != sub.entry {
			if len(snode.Preds) != len(gnode.Preds) {
				return false
			}
			for _, spred := range snode.Preds {
				found := false
				for _, gpred := range gnode.Preds {
					if m[spred.Index] == gpred.Index {
						found = true
						break
					}
				}
				if !found {
					return false
				}
			}
		}

		// Check successors.
		if snode.Index != sub.exit {
			if len(snode.Succs) != len(gnode.Succs) {
				return false
			}
			for _, ssucc := range snode.Succs {
				found := false
				for _, gsucc := range gnode.Succs {
					if m[ssucc.Index] == gsucc.Index {
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

	// Match found!
	return true
}
