// Package iso implements subgraph isomorphism search algorithms.
package iso

import (
	"fmt"
	"log"
	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/mewfork/dot"
	"github.com/mewkiz/pkg/errutil"
	"github.com/mewrev/graphs"
)

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

// Solve tries to locate a mapping from sub node name to graph node name for an
// isomorphism of sub in graph based on the given node pair candidates.
func (eq *Equation) Solve(graph *dot.Graph, sub *graphs.SubGraph) (m map[string]string, err error) {
	out := make(chan map[string]string)
	go eq.solve(graph, sub, out)
	for i := 0; i < 100; i++ {
		m := <-out
		//if m == nil {
		//	return errutil.New("unable to solve node pair equation")
		//}
		//fmt.Println("@@@ [ mapping found ] @@@@@@@@@@@@@@@@")
		if m != nil {
			spew.Dump(fmt.Sprintf("i=%d", i), m)
		} else {
			fmt.Println("<nil>")
		}
	}

	return nil, nil
}

// solve tries to locate a mapping from sub node name to graph node name for an
// isomorphism of sub in graph based on the given node pair candidates. A valid
// mapping is sent to the out channel if successful, and a nil mapping
// otherwise.
func (eq *Equation) solve(graph *dot.Graph, sub *graphs.SubGraph, out chan map[string]string) {
	for !eq.IsSolved(graph, sub) {
		// TODO: Remove debug output.
		fmt.Println("___ [ before unique pairs ] _________")
		fmt.Println()
		fmt.Println("candidate mapping:")
		fmt.Println()
		spew.Dump(eq.c)
		fmt.Println()
		fmt.Println("mapping:")
		fmt.Println()
		spew.Dump(eq.m)
		fmt.Println()

		// No candidates left.
		if len(eq.c) == 0 {
			out <- nil
			return
		}

		// Locate unique node pairs.
		ok, err := eq.SolveUnique()
		if err != nil {
			log.Println(errutil.Err(err))
			out <- nil
			return
		}
		if ok {
			continue
		}

		fmt.Println("___ [ before brute ] _________")
		fmt.Println()
		fmt.Println("candidate mapping:")
		fmt.Println()
		spew.Dump(eq.c)
		fmt.Println()
		fmt.Println("mapping:")
		fmt.Println()
		spew.Dump(eq.m)
		fmt.Println()

		// Locate the easiest node pair by brute force.
		sname, err := eq.easiest()
		if err != nil {
			log.Println(errutil.Err(err))
			out <- nil
			return
		}
		candidates := eq.c[sname]

		fmt.Println("^^^ brute ^^^")
		fmt.Println()

		// Try each node pair candidate.
		ncandidates := len(candidates)
		in := make(chan map[string]string)
		for gname := range candidates {
			fmt.Println("sname:", sname, "gname:", gname)
			fmt.Println()
			go func(eq *Equation, gname string) {
				err := eq.SetPair(sname, gname)
				if err != nil {
					log.Println(errutil.Err(err))
				}
				eq.solve(graph, sub, in)
			}(eq.Dup(), gname)
		}
		var m map[string]string
		for i := 0; i < ncandidates; i++ {
			if m != nil {
				m = <-in
			} else {
				<-in
			}
		}
		out <- m
		if m != nil {
			return
		}
	}

	out <- eq.m
}

// SolveBrute tries to solve the easiest node pair (i.e. the one with the fewest
// number of candidates) of the equation by brute force.
func (eq *Equation) SolveBrute(graph *dot.Graph, sub *graphs.SubGraph) error {
	// Locate the easiest node pair to solve.
	sname, err := eq.easiest()
	if err != nil {
		return errutil.Err(err)
	}
	candidates := eq.c[sname]

	// Try each node pair candidate.
	wg := new(sync.WaitGroup)
	wg.Add(len(candidates))
	out := make(chan map[string]string)
	for gname := range candidates {
		go brute(graph, sub, eq.Dup(), sname, gname, wg, out)
	}
	wg.Wait()

	return nil
}

func brute(graph *dot.Graph, sub *graphs.SubGraph, eq *Equation, sname, gname string, wg *sync.WaitGroup, out chan map[string]string) {
	fmt.Println("trying to solve eq with:", gname)
	err := eq.SetPair(sname, gname)
	if err != nil {
		log.Println(errutil.Err(err))
	}
	if len(eq.c) == 0 {
		out <- nil
	}
	if eq.IsSolved(graph, sub) {
		out <- eq.m
	}
	wg.Done()
}

// easiest returns the sub node name of the easiest node pair (i.e. the one with
// the fewest number of candidates) to solve.
func (eq *Equation) easiest() (string, error) {
	min := -1
	var easiest string
	for sname, candidates := range eq.c {
		if min == -1 || len(candidates) < min {
			min = len(candidates)
			easiest = sname
		}
	}
	if min < 2 {
		return "", errutil.Newf("too few candidates for brute force; expected > 2, got %d", min)
	}
	return easiest, nil
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
				return false, errutil.Err(err)
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
	if key, ok := getKey(eq.m, gname); ok {
		return errutil.Newf("invalid mapping; sub node %q and %q both map to graph node %q", key, sname, gname)
	}

	// Move unique node pair from c to m.
	eq.m[sname] = gname
	delete(eq.c, sname)

	// Remove graph node name of the unique node pair from all other node
	// pairs in c.
	for key, candidates := range eq.c {
		delete(candidates, gname)
		if len(eq.c[key]) == 0 {
			return errutil.Newf("invalid mapping; sub node %q has no candidates", key)
		}
	}

	return nil
}

// getKey returns the first key in m which maps to the value val. The boolean
// value is true if such a key could be located, and false otherwise.
func getKey(m map[string]string, val string) (key string, ok bool) {
	for key, x := range m {
		if x == val {
			return key, true
		}
	}
	return "", false
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
