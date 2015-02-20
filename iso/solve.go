package iso

import (
	"fmt"

	"github.com/mewkiz/pkg/errutil"
)

// NodePairs specifies node pair candidates and known node pairs.
type NodePairs struct {
	// mapping from sub node name to graph node name candidates.
	c map[string]map[string]bool
	// mapping from sub node name to graph node name.
	m map[string]string
}

// SolveUnique tries to locate a unique node pair in c. If successful the node
// pair is removed from c and stored in m. As the graph node name of the node
// pair is no longer a valid candidate it is removed from all other node pairs
// in c.
func (pairs *NodePairs) SolveUnique() error {
	for sname, candidates := range pairs.c {
		if len(candidates) == 1 {
			gname := pop(candidates)
			return pairs.SetPair(sname, gname)
		}
	}

	return errutil.New("unable to locate a unique node pair")

}

// SetPair marks the given node pair as known by removing it from c and storing
// it in m. As the graph node name is no longer a valid candidate it is removed
// from all other node pairs in c.
func (pairs *NodePairs) SetPair(sname, gname string) error {
	// Sanity check.
	if contains(pairs.m, gname) {
		return errutil.Newf("invalid mapping; sub node %q and %q both map to graph node %q", pairs.m[sname], sname, gname)
	}

	// Move unique node pair from c to m.
	pairs.m[sname] = gname
	delete(pairs.c, sname)

	// Remove graph node name of the unique node pair from all other node
	// pairs in c.
	for _, candidates := range pairs.c {
		delete(candidates, gname)
	}

	return nil
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
