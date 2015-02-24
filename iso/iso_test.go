package iso

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mewfork/dot"
	"github.com/mewrev/graphs"
)

func TestCandidates(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		entry     string
		want      map[string]map[string]bool
		err       string
	}{
		// i=0
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/primitives/if_else.dot",
			entry:     "A",
			want: map[string]map[string]bool{
				"A": map[string]bool{
					"A": true,
				},
				"B": map[string]bool{
					"B": true,
					"C": true,
				},
				"C": map[string]bool{
					"B": true,
					"C": true,
				},
				"D": map[string]bool{
					"D": true,
				},
			},
			err: "",
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "85",
			want: map[string]map[string]bool{
				"A": map[string]bool{
					"85": true,
				},
				"B": map[string]bool{
					"88": true,
				},
				"C": map[string]bool{
					"88": true,
				},
				"D": map[string]bool{
					"89": true,
				},
			},
			err: "",
		},
		// i=2
		{
			subPath:   "../testdata/primitives/while.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			want: map[string]map[string]bool{
				"A": map[string]bool{
					"71": true,
				},
				"B": map[string]bool{
					"74": true,
				},
				"C": map[string]bool{
					"74": true,
				},
			},
			err: "",
		},
		// i=3
		{
			subPath:   "../testdata/primitives/while.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "89",
			want: map[string]map[string]bool{
				"A": map[string]bool{
					"89": true,
				},
				"B": map[string]bool{
					"92": true,
					"93": true,
				},
				"C": map[string]bool{
					"92": true,
					"93": true,
				},
			},
			err: "",
		},
		// i=4
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			want: map[string]map[string]bool{
				"A": map[string]bool{
					"71": true,
				},
				"B": map[string]bool{
					"74": true,
				},
				"C": map[string]bool{
					"75": true,
				},
			},
			err: "",
		},
		// i=5
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "foo",
			want:      nil,
			err:       `unable to locate entry node "foo" in graph`,
		},
		// i=6
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "97",
			want:      nil,
			err:       `invalid entry node candidate "97"; expected 2 successors, got 1`,
		},
		// i=7
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "68",
			want:      nil,
			err:       "incomplete candidate mapping; expected 4 map entites, got 1",
		},
	}

	for i, g := range golden {
		sub, err := graphs.ParseSubGraph(g.subPath)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		graph, err := dot.ParseFile(g.graphPath)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		eq, err := Candidates(graph, g.entry, sub)
		if !sameError(err, g.err) {
			t.Errorf("i=%d: error mismatch; expected %v, got %v", i, g.err, err)
			continue
		} else if err != nil {
			// Expected error, check next test case.
			continue
		}
		if !reflect.DeepEqual(eq.c, g.want) {
			t.Errorf("i=%d: candidate map mismatch; expected %v, got %v", i, g.want, eq.c)
		}
	}
}

func TestSolve(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		entry     string
		wants     []map[string]string
	}{
		// i=0
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/primitives/if_else.dot",
			entry:     "A",
			wants: []map[string]string{
				{
					"A": "A",
					"B": "B",
					"C": "C",
					"D": "D",
				},
				{
					"A": "A",
					"B": "C",
					"C": "B",
					"D": "D",
				},
			},
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "85",
			wants:     []map[string]string{nil},
		},
		// i=2
		{
			subPath:   "../testdata/primitives/while.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			wants:     []map[string]string{nil},
		},
		// i=3
		{
			subPath:   "../testdata/primitives/while.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "89",
			wants: []map[string]string{
				{
					"A": "89",
					"B": "92",
					"C": "93",
				},
			},
		},
		// i=4
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			wants: []map[string]string{
				{
					"A": "71",
					"B": "74",
					"C": "75",
				},
			},
		},
	}

loop:
	for i, g := range golden {
		graph, err := dot.ParseFile(g.graphPath)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		sub, err := graphs.ParseSubGraph(g.subPath)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		eq, err := Candidates(graph, g.entry, sub)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		m, err := eq.Solve(graph, sub)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		for _, want := range g.wants {
			if reflect.DeepEqual(m, want) {
				continue loop
			}
		}
		t.Errorf("i=%d: node pair map mismatch; expected one of %v, got %v", i, g.wants, m)
	}
}

func TestEquationSetPair(t *testing.T) {
	golden := []struct {
		in           *Equation
		sname, gname string
		want         *Equation
		err          string
	}{
		// i=0
		{
			in: &Equation{
				c: map[string]map[string]bool{
					"A": map[string]bool{
						"A": true,
					},
					"B": map[string]bool{
						"B": true,
						"C": true,
					},
					"C": map[string]bool{
						"B": true,
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{},
			},
			sname: "A", gname: "A",
			want: &Equation{
				c: map[string]map[string]bool{
					"B": map[string]bool{
						"B": true,
						"C": true,
					},
					"C": map[string]bool{
						"B": true,
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			err: "",
		},
		// i=1
		{
			in: &Equation{
				c: map[string]map[string]bool{
					"B": map[string]bool{
						"B": true,
						"C": true,
					},
					"C": map[string]bool{
						"B": true,
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			sname: "B", gname: "B",
			want: &Equation{
				c: map[string]map[string]bool{
					"C": map[string]bool{
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"B": "B",
				},
			},
			err: "",
		},
		// i=2
		{
			in: &Equation{
				c: map[string]map[string]bool{
					"B": map[string]bool{
						"B": true,
						"C": true,
					},
					"C": map[string]bool{
						"B": true,
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			sname: "B", gname: "C",
			want: &Equation{
				c: map[string]map[string]bool{
					"C": map[string]bool{
						"B": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"B": "C",
				},
			},
			err: "",
		},
		// i=3
		{
			in: &Equation{
				c: map[string]map[string]bool{
					"A": map[string]bool{
						"A": true,
						"D": true,
					},
					"C": map[string]bool{
						"A": true,
						"C": true,
					},
					"D": map[string]bool{
						"D": true,
					},
				},
				m: map[string]string{},
			},
			sname: "A", gname: "D",
			want: nil,
			err:  `invalid mapping; sub node "D" has no candidates`,
		},
		// i=3
		{
			in: &Equation{
				c: map[string]map[string]bool{
					"A": map[string]bool{
						"0": true,
						"1": true,
					},
				},
				m: map[string]string{
					"B": "1",
				},
			},
			sname: "A", gname: "1",
			want: nil,
			err:  `invalid mapping; sub node "B" and "A" both map to graph node "1"`,
		},
	}

	for i, g := range golden {
		err := g.in.SetPair(g.sname, g.gname)
		if !sameError(err, g.err) {
			t.Errorf("i=%d: error mismatch; expected %v, got %v", i, g.err, err)
			continue
		} else if err != nil {
			// Expected error, check next test case.
			continue
		}
		if !reflect.DeepEqual(g.want, g.in) {
			t.Errorf("i=%d: node pair equation mismatch; expected %v, got %v", i, g.want, g.in)
		}
	}
}

// sameError returns true if err is represented by the string s, and false
// otherwise. Some error messages contains "file:line" prefixes and suffixes
// from external functions, e.g.
//
//    github.com/mewrev/graphs/iso.Candidates (solve.go:53): error: unable to locate entry node "foo" in graph
//    unable to parse integer constant "foo"; strconv.ParseInt: parsing "foo": invalid syntax`
//
// For this reason s matches the error if it is a non-empty substring of err.
func sameError(err error, s string) bool {
	t := ""
	if err != nil {
		if len(s) == 0 {
			return false
		}
		t = err.Error()
	}
	return strings.Contains(t, s)
}
