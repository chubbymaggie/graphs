package iso

import "github.com/mewkiz/pkg/errutil"

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
