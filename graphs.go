// Package graphs implements subgraph isomorphism search algorithms.
package graphs

import (
	"log"

	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
)

// SubGraph represents a subgraph with a dedicated entry and exit node. Incoming
// edges to entry and outgoing edges from exit are ignored when searching for
// isomorphisms of the subgraph.
type SubGraph struct {
	*dot.Graph
	entry, exit string
}

// ParseSubGraph parses the provided DOT file into a subgraph with a dedicated
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
func ParseSubGraph(path string) (*SubGraph, error) {
	graph, err := dot.ParseFile(path)
	if err != nil {
		return nil, err
	}
	return NewSubGraph(graph)
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
				return nil, errutil.Newf(`redefinition of node with "entry" label; previous node %q, new node %q`, sub.entry, node.Name)
			}
			sub.entry = node.Name
			hasEntry = true
		case "exit":
			if hasExit {
				return nil, errutil.Newf(`redefinition of node with "exit" label; previous node %q, new node %q`, sub.exit, node.Name)
			}
			sub.exit = node.Name
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

// Entry returns the entry node name in the subgraph.
func (sub *SubGraph) Entry() string {
	return sub.entry
}

// Exit returns the exit node name in the subgraph.
func (sub *SubGraph) Exit() string {
	return sub.exit
}

// Search tries to locate an isomorphism of sub in graph. If successful it
// returns the mapping from sub node name to graph node name of the first
// isomorphism located. The boolean value is true if such a mapping could be
// located, and false otherwise.
func Search(graph *dot.Graph, sub *SubGraph) (m map[string]string, ok bool) {
	for _, node := range graph.Nodes.Nodes {
		m, ok = Isomorphism(graph, node.Name, sub)
		if ok {
			return m, true
		}
	}
	return nil, false
}

// Isomorphism returns a mapping from sub node name to graph node name if there
// exists an isomorphism of sub in graph (starting at the entry node). The
// boolean value is true if such a mapping could be located, and false
// otherwise.
func Isomorphism(graph *dot.Graph, entry string, sub *SubGraph) (m map[string]string, ok bool) {
	m = make(map[string]string)
	visited := make(map[pair]bool)
	g, ok := graph.Nodes.Lookup[entry]
	if !ok {
		log.Printf("graphs.Isomorphism: unable to locate entry node %q in graph.\n", entry)
		return nil, false
	}
	s, ok := sub.Nodes.Lookup[sub.entry]
	if !ok {
		log.Printf("graphs.Isomorphism: unable to locate entry node %q in subgraph.\n", sub.entry)
		return nil, false
	}
	if isIsomorphism(g, s, graph, sub, m, visited) {
		return m, true
	}
	return nil, false
}

// pair is a key-value pair.
type pair struct {
	key, val string
}

// isIsomorphism returns true if g is an isomorphism of s, where g is a node of
// graph, s is a node of sub and m is a mapping from sub node name to graph node
// name. Incoming edges to entry and outgoing edges from exit are ignored when
// searching for isomorphisms of sub.
func isIsomorphism(g, s *dot.Node, graph *dot.Graph, sub *SubGraph, m map[string]string, visited map[pair]bool) bool {
	// TODO: Check for loops?
	// TODO: Check for duplicate val in m and only add if not already present.
	// TODO: Take edge labels (e.g. conditional branches) into account?

	// HACK: the visited map is used to prevent cycles. Find a cleaner solution
	// and remove visited entirely!

	// Create mapping from sub node name to graph node name by trying possible
	// candidates.
	if len(m) != len(sub.Nodes.Nodes) {
		if s.Name == sub.entry {
			// Add mapping for entry node.
			if _, ok := m[s.Name]; !ok {
				// TODO: How are graphs with entry points having outgoing edges to
				// themselves handled?
				m[s.Name] = g.Name
			}
		} else {
			// Verify predecessor count.
			if len(s.Preds) != len(g.Preds) {
				return false
			}
		}

		// Add candidate mapping for successors.
		if s.Name != sub.exit {
			if len(s.Succs) != len(g.Succs) {
				return false
			}
			for _, ssucc := range s.Succs {
				for _, gsucc := range g.Succs {
					if _, ok := m[ssucc.Name]; !ok {
						m[ssucc.Name] = gsucc.Name
					}
					pair := pair{key: ssucc.Name, val: gsucc.Name}
					if visited[pair] {
						continue
					}
					visited[pair] = true
					if isIsomorphism(gsucc, ssucc, graph, sub, m, visited) {
						return true
					}
					delete(m, ssucc.Name)
				}
			}
		}
	}

	// Complete mapping?
	if len(m) != len(sub.Nodes.Nodes) {
		return false
	}

	// Check if m correctly maps the nodes of sub onto the nodes of graph.
	for sname, gname := range m {
		snode, ok := sub.Nodes.Lookup[sname]
		if !ok {
			log.Printf("graphs.isIsomorphism: unable to locate node %q in subgraph.\n", sname)
			return false
		}
		gnode, ok := graph.Nodes.Lookup[gname]
		if !ok {
			log.Printf("graphs.isIsomorphism: unable to locate node %q in graph.\n", gname)
			return false
		}

		// Check predecessors.
		if snode.Name != sub.entry {
			if len(snode.Preds) != len(gnode.Preds) {
				return false
			}
			for _, spred := range snode.Preds {
				found := false
				for _, gpred := range gnode.Preds {
					if m[spred.Name] == gpred.Name {
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
		if snode.Name != sub.exit {
			if len(snode.Succs) != len(gnode.Succs) {
				return false
			}
			for _, ssucc := range snode.Succs {
				found := false
				for _, gsucc := range gnode.Succs {
					if m[ssucc.Name] == gsucc.Name {
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
