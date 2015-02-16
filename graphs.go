// Package graphs implements subgraph isomorphism search algorithms.
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
	gnodes, snodes := graph.Nodes.Nodes, sub.Nodes.Nodes
	for _, node := range gnodes {
		label, ok := node.Attrs["label"]
		if !ok {
			continue
		}
		switch label {
		case "entry":
			if hasEntry {
				return nil, errutil.Newf(`redefinition of node with "entry" label; previous node %q, new node %q`, snodes[sub.entry].Name, gnodes[node.Index].Name)
			}
			sub.entry = node.Index
			hasEntry = true
		case "exit":
			if hasExit {
				return nil, errutil.Newf(`redefinition of node with "exit" label; previous node %q, new node %q`, snodes[sub.exit].Name, gnodes[node.Index].Name)
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
	visited := make(map[pair]bool)
	g := graph.Nodes.Nodes[entry]
	s := sub.Nodes.Nodes[sub.entry]
	if isIsomorphism(g, s, graph, sub, m, visited) {
		return m, true
	}
	return nil, false
}

// pair is a key-value pair.
type pair struct {
	key, val int
}

// isIsomorphism returns true if g is an isomorphism of s, where g is a node of
// graph, s is a node of sub and m is a mapping from sub node index to graph
// node index. Incoming edges to entry and outgoing edges from exit are ignored
// when searching for isomorphisms of sub.
func isIsomorphism(g, s *dot.Node, graph *dot.Graph, sub *SubGraph, m map[int]int, visited map[pair]bool) bool {
	// TODO: Check for loops?
	// TODO: Check for duplicate val in m and only add if not already present.
	// TODO: Take edge labels (e.g. conditional branches) into account?

	// HACK: the visited map is used to prevent cycles. Find a cleaner solution
	// and remove visited entirely!

	// Create mapping from sub node index to graph node index by trying possible
	// candidates.
	gnodes, snodes := graph.Nodes.Nodes, sub.Nodes.Nodes
	if len(m) != len(snodes) {
		if s.Index == sub.entry {
			// Add mapping for entry node.
			if _, ok := m[s.Index]; !ok {
				// TODO: How are graphs with entry points having outgoing edges to
				// themselves handled?
				m[s.Index] = g.Index
			}
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
				for _, gsucc := range g.Succs {
					if _, ok := m[ssucc.Index]; !ok {
						m[ssucc.Index] = gsucc.Index
					}
					pair := pair{key: ssucc.Index, val: gsucc.Index}
					if visited[pair] {
						continue
					}
					visited[pair] = true
					if isIsomorphism(gsucc, ssucc, graph, sub, m, visited) {
						return true
					}
					delete(m, ssucc.Index)
				}
			}
		}
	}

	// Complete mapping?
	if len(m) != len(snodes) {
		return false
	}

	// Check if m correctly maps the nodes of sub onto the nodes of graph.
	for sidx, gidx := range m {
		snode, gnode := snodes[sidx], gnodes[gidx]

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
