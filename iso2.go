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
	c := make(map[string]map[string]bool)
	locate(g, s, sub, c)
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

	// Prevent infinite cycles.
	if _, ok := c[s.Name]; ok {
		if c[s.Name][g.Name] {
			log.Printf("already visited (%q=%q)\n", s.Name, g.Name)
			return
		}
	}

	// Add node pair candidate. Add entry node pair exactly once.
	if _, ok := c[s.Name]; !ok {
		c[s.Name] = map[string]bool{
			g.Name: true,
		}
	} else if s.Name != sub.entry {
		c[s.Name][g.Name] = true
	}

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

// Solve returns a mapping from sub node name to graph node name for an
// isomorphism of sub in graph based on the given node pair candidates.
func Solve(graph *dot.Graph, sub *SubGraph, c map[string]map[string]bool) (map[string]string, error) {
	// Sanity check.
	if len(c) != len(sub.Nodes.Nodes) {
		return nil, errutil.Newf("incomplete mapping; expected %d map entities, got %d", len(sub.Nodes.Nodes), len(c))
	}

	m := make(map[string]string)
	for {
		// Locate unique node pairs.
		err := solveUniqPair(c, m)
		if err != nil {
			if len(c) > 0 {
				fmt.Println("~~~ [ map ] ~~~")
				spew.Dump(m)
				fmt.Println("~~~ [ needs attention ] ~~~")
				spew.Dump(c)
				panic("foo")
			}
			return nil, errutil.Err(err)
		}

		if valid(graph, sub, m) {
			return m, nil
		}
	}
}

// solveUniqPair tries to locate a unique node pair in c. If successful the node
// pair is removed from c and stored in m. As the graph node name of the node
// pair is no longer a valid candidate it is removed from all other node pairs
// in c.
func solveUniqPair(c map[string]map[string]bool, m map[string]string) error {
	for sname, candidates := range c {
		if len(candidates) != 1 {
			continue
		}

		gname := pop(candidates)
		if contains(m, gname) {
			return errutil.Newf("invalid mapping; sub node %q and %q both map to graph node %q", m[sname], sname, gname)
		}

		// Move unique node pair from c to m.
		m[sname] = gname
		delete(c, sname)

		// Remove graph node name of the unique node pair from all other node
		// pairs in c.
		for _, candidates := range c {
			delete(candidates, gname)
		}

		return nil
	}

	return errutil.New("unable to locate a unique node pair")
}

// contains returns true if m contains the value val, and false otherwise.
func contains(m map[string]string, val string) bool {
	for _, x := range m {
		if x == val {
			return true
		}
	}
	return false
}

// pop returns the only key in m.
func pop(m map[string]bool) string {
	if len(m) != 1 {
		panic(fmt.Sprintf("invalid map length; expected 1, got %d", len(m)))
	}
	for key := range m {
		return key
	}
	panic("unreachable")
}
