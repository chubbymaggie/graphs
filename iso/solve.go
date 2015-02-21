// Package iso implements subgraph isomorphism search algorithms.
package iso

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewrev/graphs"
)

// Equation specifies an equation of node pair candidates and known node pairs.
type Equation struct {
	// mapping from sub node name to graph node name candidates.
	c map[string]map[string]bool
	// mapping from sub node name to graph node name.
	m map[string]string
}

// Dup returns a copy of eq.
func (eq *Equation) Dup() *Equation {
	// Duplicate node pair candidates.
	c := make(map[string]map[string]bool)
	for sname, candidates := range eq.c {
		c[sname] = make(map[string]bool)
		for gname, val := range candidates {
			c[sname][gname] = val
		}
	}

	// Duplicate node pairs.
	m := make(map[string]string)
	for sname, gname := range eq.m {
		m[sname] = gname
	}

	return &Equation{c: c, m: m}
}

// TODO: Remove the C and M methods.
func (eq *Equation) C() map[string]map[string]bool { return eq.c }
func (eq *Equation) M() map[string]string          { return eq.m }

// Candidates locates node pair candidates for an isomorphism of sub in graph
// which starts at the entry node.
func Candidates(graph *dot.Graph, entry string, sub *graphs.SubGraph) (*Equation, error) {
	// Sanity checks.
	g, ok := graph.Nodes.Lookup[entry]
	if !ok {
		return nil, errutil.Newf("unable to locate entry node %q in graph", entry)
	}
	s, ok := sub.Nodes.Lookup[sub.Entry()]
	if !ok {
		return nil, errutil.Newf("unable to locate entry node %q in sub", sub.Entry())
	}
	if !isPotential(g, s, sub) {
		return nil, errutil.Newf("invalid entry node candidate %q; expected %d successors, got %d", g.Name, len(s.Succs), len(g.Succs))
	}

	// Locate candidate node pairs.
	eq := &Equation{
		c: make(map[string]map[string]bool),
		m: make(map[string]string),
	}
	eq.findCandidates(g, s, sub)
	if len(eq.c) != len(sub.Nodes.Nodes) {
		fmt.Println("### [ incomplete mapping ] ###")
		spew.Dump(eq.c)
		fmt.Println("### [/ incomplete mapping ] ###")
		return nil, errutil.Newf("incomplete candidate mapping; expected %d map entites, got %d", len(sub.Nodes.Nodes), len(eq.c))
	}

	return eq, nil
}

// findCandidates recursively locates potential node pairs (g and s) for an
// isomorphism of sub in graph and adds them to c.
func (eq *Equation) findCandidates(g, s *dot.Node, sub *graphs.SubGraph) {
	// Exit early for impossible node pairs.
	if !isPotential(g, s, sub) {
		return
	}

	// Prevent infinite cycles.
	if _, ok := eq.c[s.Name]; ok {
		if eq.c[s.Name][g.Name] {
			// TODO: Remove debug output.
			log.Printf("already visited (%q=%q)\n", s.Name, g.Name)
			return
		}
	}

	// Add node pair candidate. Add entry node pair exactly once.
	if _, ok := eq.c[s.Name]; !ok {
		eq.c[s.Name] = map[string]bool{
			g.Name: true,
		}
	} else if s.Name != sub.Entry() {
		eq.c[s.Name][g.Name] = true
	}

	// Recursively locate candidate successor pairs.
	for _, ssucc := range s.Succs {
		for _, gsucc := range g.Succs {
			eq.findCandidates(gsucc, ssucc, sub)
		}
	}
}

// Solve tries to locate a mapping from sub node name to graph node name for an
// isomorphism of sub in graph based on the given node pair candidates.
func (eq *Equation) Solve(graph *dot.Graph, sub *graphs.SubGraph) error {
	// Sanity check.
	if len(eq.c) != len(sub.Nodes.Nodes) {
		return errutil.Newf("incomplete candidate mapping; expected %d map entites, got %d", len(sub.Nodes.Nodes), len(eq.c))
	}

	for {
		// Locate unique node pairs.
		ok, err := eq.SolveUnique()
		if err != nil {
			return errutil.Err(err)
		}
		if ok {
			continue
		}

		if len(eq.c) > 0 {
			fmt.Println("~~~ [ map ] ~~~")
			spew.Dump(eq.m)
			fmt.Println("~~~ [ needs attention ] ~~~")
			spew.Dump(eq.c)
			panic("foo")
		}

		if eq.IsSolved(graph, sub) {
			return nil
		}
	}
}

// SolveUnique tries to locate a unique node pair in c. If successful the node
// pair is removed from c and stored in m. As the graph node name of the node
// pair is no longer a valid candidate it is removed from all other node pairs
// in c.
func (eq *Equation) SolveUnique() (ok bool, err error) {
	for sname, candidates := range eq.c {
		if len(candidates) == 1 {
			gname := pop(candidates)
			err := eq.SetPair(sname, gname)
			if err != nil {
				return false, err
			}
			return true, nil
		}
	}

	return false, nil
}

// SetPair marks the given node pair as known by removing it from c and storing
// it in m. As the graph node name is no longer a valid candidate it is removed
// from all other node pairs in c.
func (eq *Equation) SetPair(sname, gname string) error {
	// Sanity check.
	if contains(eq.m, gname) {
		return errutil.Newf("invalid mapping; sub node %q and %q both map to graph node %q", eq.m[sname], sname, gname)
	}

	// Move unique node pair from c to m.
	eq.m[sname] = gname
	delete(eq.c, sname)

	// Remove graph node name of the unique node pair from all other node
	// pairs in c.
	for _, candidates := range eq.c {
		delete(candidates, gname)
	}

	return nil
}

// isPotential returns true if the graph node g is a potential candidate for the
// sub node s, and false otherwise.
func isPotential(g, s *dot.Node, sub *graphs.SubGraph) bool {
	// Verify predecessors.
	if s.Name != sub.Entry() && len(g.Preds) != len(s.Preds) {
		return false
	}

	// Verify successors.
	if s.Name != sub.Exit() && len(g.Succs) != len(s.Succs) {
		return false
	}

	return true
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
