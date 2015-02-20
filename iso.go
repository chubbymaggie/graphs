package graphs

import (
	"fmt"
	"log"
	"sort"

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

	// Check for duplicate values.
	if hasDup(m) {
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
			if len(s.Preds) != len(g.Preds) {
				return false
			}
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
			if len(s.Succs) != len(g.Succs) {
				return false
			}
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

// hasDup returns true if m contains a duplicate value.
func hasDup(m map[string]string) bool {
	vals := make(map[string]bool, len(m))
	for _, v := range m {
		if vals[v] {
			return true
		}
		vals[v] = true
	}
	return false
}

// Isomorphism returns a mapping from sub node name to graph node name if there
// exists an isomorphism of sub in graph which starts at the entry node. The
// boolean value is true if such a mapping could be located, and false
// otherwise.
func Isomorphism(graph *dot.Graph, entry string, sub *SubGraph) (m map[string]string, ok bool) {
	g, ok := graph.Nodes.Lookup[entry]
	if !ok {
		err := errutil.Newf("unable to locate entry node %q in graph", entry)
		log.Println(err)
		return nil, false
	}
	s, ok := sub.Nodes.Lookup[sub.entry]
	if !ok {
		err := errutil.Newf("unable to locate entry node %q in sub", sub.entry)
		log.Println(err)
		return nil, false
	}
	c := make(map[string]string)
	visited := make(map[string]bool)
	find(g, s, graph, sub, c, &m, visited)
	if m != nil {
		return m, true
	}
	return nil, false
}

// pair is a key-value node pair.
type pair struct {
	// sub node name.
	key string
	// graph node name.
	val string
}

// find searches for an isomorphism of sub in graph, where g and s represent a
// candidate node pair from graph and sub respectively, c a candidate mapping
// and out the output for a valid mapping. To prevent cycles visited keeps track
// of already visited node pairs.
func find(g, s *dot.Node, graph *dot.Graph, sub *SubGraph, c map[string]string, out *map[string]string, visited map[string]bool) {
	// Valid mapping already located?
	if *out != nil {
		return
	}

	if s.Name == sub.entry {
		if _, ok := c[s.Name]; !ok {
			c[s.Name] = g.Name
		}
	} else {
		c[s.Name] = g.Name
	}

	fmt.Println("candidate mapping:", c)

	// Node pair already visited?
	if visited[enc(c)] {
		log.Println("already visited.")
		return
	}

	// Validate candidate mapping.
	if len(c) == len(sub.Nodes.Nodes) {
		if valid(graph, sub, c) {
			*out = c
			return
		}
	}

	// Verify predecessors.
	if s.Name != sub.entry {
		if len(s.Preds) != len(g.Preds) {
			return
		}
	}

	// Verify successors.
	if s.Name != sub.exit {
		if len(s.Succs) != len(g.Succs) {
			return
		}
		for _, ssucc := range sortNodes(s.Succs) {
			for _, gsucc := range sortNodes(g.Succs) {
				find(gsucc, ssucc, graph, sub, c, out, visited)
			}
		}
	}
}

func enc(m map[string]string) string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	s := ""
	for _, key := range keys {
		s += fmt.Sprintf("%q=%q, ", key, m[key])
	}
	return s
}
