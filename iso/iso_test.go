package iso

import (
	"reflect"
	"strings"
	"testing"

	"decomp.org/x/graphs"
	"github.com/mewfork/dot"
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
				"A": {
					"A": true,
				},
				"B": {
					"B": true,
					"C": true,
				},
				"C": {
					"B": true,
					"C": true,
				},
				"D": {
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
				"A": {
					"85": true,
				},
				"B": {
					"88": true,
				},
				"C": {
					"88": true,
				},
				"D": {
					"89": true,
				},
			},
			err: "",
		},
		// i=2
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			want: map[string]map[string]bool{
				"A": {
					"71": true,
				},
				"B": {
					"74": true,
				},
				"C": {
					"74": true,
				},
			},
			err: "",
		},
		// i=3
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "89",
			want: map[string]map[string]bool{
				"A": {
					"89": true,
				},
				"B": {
					"92": true,
					"93": true,
				},
				"C": {
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
				"A": {
					"71": true,
				},
				"B": {
					"74": true,
				},
				"C": {
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
		eq, err := candidates(graph, g.entry, sub)
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

func TestSolveBrute(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		entry     string
		wants     []map[string]string
		err       string
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
			err: "",
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "85",
			wants:     []map[string]string{nil},
			err:       "unable to locate node pair mapping",
		},
		// i=2
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			wants:     []map[string]string{nil},
			err:       "unable to locate node pair mapping",
		},
		// i=3
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
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
		eq, err := candidates(graph, g.entry, sub)
		if err != nil {
			t.Errorf("i=%d: %v", i, err)
			continue
		}
		m, err := eq.solveBrute(graph, sub)
		if !sameError(err, g.err) {
			t.Errorf("i=%d: error mismatch; expected %v, got %v", i, g.err, err)
			continue
		} else if err != nil {
			// Expected error, check next test case.
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
		in           *equation
		sname, gname string
		want         *equation
		err          string
	}{
		// i=0
		{
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
						"A": true,
					},
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
				},
				m: map[string]string{},
			},
			sname: "A", gname: "A",
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
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
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			sname: "B", gname: "B",
			want: &equation{
				c: map[string]map[string]bool{
					"C": {
						"C": true,
					},
					"D": {
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
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			sname: "B", gname: "C",
			want: &equation{
				c: map[string]map[string]bool{
					"C": {
						"B": true,
					},
					"D": {
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
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
						"A": true,
						"D": true,
					},
					"C": {
						"A": true,
						"C": true,
					},
					"D": {
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
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
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
		err := g.in.setPair(g.sname, g.gname)
		if !sameError(err, g.err) {
			t.Errorf("i=%d: error mismatch; expected %v, got %v", i, g.err, err)
			continue
		} else if err != nil {
			// Expected error, check next test case.
			continue
		}
		if !reflect.DeepEqual(g.in, g.want) {
			t.Errorf("i=%d: node pair equation mismatch; expected %v, got %v", i, g.want, g.in)
		}
	}
}

func TestEquationDup(t *testing.T) {
	golden := []struct {
		in         *equation
		ckey, mkey string
		want       *equation
	}{
		// i=0
		{
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
						"A": true,
					},
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"D": "D",
					"E": "E",
				},
			},
			ckey: "A", mkey: "D",
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"E": "E",
				},
			},
		},
		// i=1
		{
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
					"E": "E",
				},
			},
			ckey: "B", mkey: "E",
			want: &equation{
				c: map[string]map[string]bool{
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
				},
			},
		},
	}

	for i, g := range golden {
		got := g.in.dup()
		if !reflect.DeepEqual(got, g.in) {
			t.Errorf("i=%d: equation copy differs from original; expected %v, got %v", i, g.in, got)
			continue
		}
		delete(got.c, g.ckey)
		delete(got.m, g.mkey)
		if reflect.DeepEqual(got.c, g.in.c) {
			t.Errorf("i=%d: copy refers to the same candidate node pair map as the original equation", i)
		}
		if reflect.DeepEqual(got.m, g.in.m) {
			t.Errorf("i=%d: copy refers to the same known node pair map as the original equation", i)
		}
		if !reflect.DeepEqual(got, g.want) {
			t.Errorf("i=%d: unable to delete keys from equation copy", i)
		}
	}
}

func TestEquationSolveUnique(t *testing.T) {
	golden := []struct {
		in   *equation
		want *equation
		ok   bool
		err  string
	}{
		// i=0
		{
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
						"A": true,
					},
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"D": "D",
					"E": "E",
				},
			},
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
					"E": "E",
				},
			},
			ok:  true,
			err: "",
		},
		// i=1
		{
			in: &equation{
				c: map[string]map[string]bool{
					"A": {
						"A": true,
					},
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
					"E": {
						"E": true,
					},
				},
				m: map[string]string{},
			},
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
					"E": {
						"E": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			ok:  true,
			err: "",
		},
		// i=2
		{
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"D": {
						"D": true,
					},
					"E": {
						"E": true,
					},
				},
				m: map[string]string{
					"A": "A",
				},
			},
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"E": {
						"E": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
				},
			},
			ok:  true,
			err: "",
		},
		// i=3
		{
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
					"E": {
						"E": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
				},
			},
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
					"E": "E",
				},
			},
			ok:  true,
			err: "",
		},
		// i=4
		{
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
					"E": "E",
				},
			},
			want: &equation{
				c: map[string]map[string]bool{
					"B": {
						"B": true,
						"C": true,
					},
					"C": {
						"B": true,
						"C": true,
					},
				},
				m: map[string]string{
					"A": "A",
					"D": "D",
					"E": "E",
				},
			},
			ok:  false,
			err: "",
		},
		// i=5
		{
			in: &equation{
				c: map[string]map[string]bool{
					"B": {
						"0": true,
					},
					"C": {
						"1": true,
						"2": true,
					},
				},
				m: map[string]string{
					"A": "0",
					"D": "3",
					"E": "4",
				},
			},
			want: nil,
			ok:   false,
			err:  `invalid mapping; sub node "A" and "B" both map to graph node "0"`,
		},
	}

	for i, g := range golden {
		ok, err := g.in.solveUnique()
		if !sameError(err, g.err) {
			t.Errorf("i=%d: error mismatch; expected %v, got %v", i, g.err, err)
			continue
		} else if err != nil {
			// Expected error, check next test case.
			continue
		}
		if ok != g.ok {
			t.Errorf("i=%d: ok mismatch; expected %v, got %v", i, g.ok, ok)
			continue
		}
		if !reflect.DeepEqual(g.in, g.want) {
			t.Errorf("i=%d: node pair equation mismatch; expected %v, got %v", i, g.want, g.in)
		}
	}
}

func TestEquationIsValid(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		eq        *equation
		want      bool
	}{
		// i=0
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "71",
					"B": "74",
					"C": "75",
				},
			},
			want: true,
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "17",
					"B": "24",
					"C": "32",
				},
			},
			want: true,
		},
		// i=2
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "89",
					"B": "92",
					"C": "93",
				},
			},
			want: false,
		},
		// i=3
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "94",
					"B": "97",
					"C": "98",
				},
			},
			want: false,
		},
		// i=4
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "282",
					"B": "292",
					"C": "287",
					"D": "299",
				},
			},
			want: true,
		},
		// i=5
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "282",
					"B": "287",
					"C": "292",
					"D": "299",
				},
			},
			want: true,
		},
		// i=6
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			eq: &equation{
				m: map[string]string{
					"A": "438",
					"B": "446",
					"C": "443",
					"D": "447",
				},
			},
			want: true,
		},
		// i=7
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			eq: &equation{
				m: map[string]string{
					"A": "438",
					"B": "443",
					"C": "446",
					"D": "447",
				},
			},
			want: true,
		},
		// i=8
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			eq: &equation{
				m: map[string]string{
					"A": "487",
					"B": "492",
					"C": "495",
					"D": "496",
				},
			},
			want: true,
		},
		// i=9
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			eq: &equation{
				m: map[string]string{
					"A": "487",
					"B": "495",
					"C": "492",
					"D": "496",
				},
			},
			want: true,
		},
		// i=10
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			eq: &equation{
				m: map[string]string{
					"A": "124",
					"B": "134",
					"C": "126",
					"D": "145",
				},
			},
			want: false,
		},
		// i=11
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			eq: &equation{
				m: map[string]string{
					"A": "740",
					"B": "760",
				},
			},
			want: true,
		},
		// i=12
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			eq: &equation{
				m: map[string]string{
					"A": "761",
					"B": "762",
				},
			},
			want: false,
		},
		// i=13
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "191",
					"B": "194",
					"C": "196",
				},
			},
			want: true,
		},
		// i=14
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "370",
					"B": "378",
					"C": "374",
				},
			},
			want: false,
		},
		// i=15
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "526",
					"B": "530",
					"C": "539",
				},
			},
			want: false,
		},
		// i=16
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "611",
					"B": "615",
					"C": "615",
					"D": "631",
				},
			},
			want: false,
		},
		// i=17
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			eq: &equation{
				m: map[string]string{
					"A": "611",
					"B": "615",
					"D": "631",
				},
			},
			want: false,
		},
		// i=18
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			eq: &equation{
				m: map[string]string{
					"A": "20",
					"B": "25",
					"C": "34",
				},
			},
			want: false,
		},
		// i=19
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "39", // 48
					"B": "44",
					"C": "52", // 45
				},
			},
			want: false,
		},
		// i=20
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			eq: &equation{
				m: map[string]string{
					"A": "39",
					"B": "44",
					"C": "45",
				},
			},
			want: false,
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
		got := g.eq.isValid(graph, sub)
		if got != g.want {
			t.Errorf("i=%d: ok mismatch; expected %v, got %v", i, g.want, got)
			continue
		}
	}
}

func TestIsomorphism(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		entry     string
		m         map[string]string
		ok        bool
	}{
		// i=0
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "71",
			m: map[string]string{
				"A": "71",
				"B": "74",
				"C": "75",
			},
			ok: true,
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "17",
			m: map[string]string{
				"A": "17",
				"B": "24",
				"C": "32",
			},
			ok: true,
		},
		// i=2
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "89",
			m:         nil,
			ok:        false,
		},
		// i=3
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "94",
			m:         nil,
			ok:        false,
		},
		// i=4
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "282",
			m: map[string]string{
				"A": "282",
				"B": "287",
				"C": "292",
				"D": "299",
			},
			ok: true,
		},
		// i=5
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			entry:     "438",
			m: map[string]string{
				"A": "438",
				"B": "443",
				"C": "446",
				"D": "447",
			},
			ok: true,
		},
		// i=6
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			entry:     "487",
			m: map[string]string{
				"A": "487",
				"B": "492",
				"C": "495",
				"D": "496",
			},
			ok: true,
		},
		// i=7
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			entry:     "124",
			m:         nil,
			ok:        false,
		},
		// i=8
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			entry:     "740",
			m: map[string]string{
				"A": "740",
				"B": "760",
			},
			ok: true,
		},
		// i=9
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			entry:     "761",
			m:         nil,
			ok:        false,
		},
		// i=10
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "191",
			m: map[string]string{
				"A": "191",
				"B": "194",
				"C": "196",
			},
			ok: true,
		},
		// i=11
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "370",
			m:         nil,
			ok:        false,
		},
		// i=12
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "526",
			m:         nil,
			ok:        false,
		},
		// i=13
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "611",
			m:         nil,
			ok:        false,
		},
		// i=14
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			entry:     "611",
			m:         nil,
			ok:        false,
		},
		// i=15
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			entry:     "20",
			m:         nil,
			ok:        false,
		},
		// i=16
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			entry:     "39",
			m:         nil,
			ok:        false,
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
		m, ok := Isomorphism(graph, g.entry, sub)
		if ok != g.ok {
			t.Errorf("i=%d: ok mismatch; expected %v, got %v", i, g.ok, ok)
			continue
		}
		if !reflect.DeepEqual(m, g.m) {
			t.Errorf("i=%d: node pair mapping mismatch; expected %v, got %v", i, g.m, m)
		}
	}
}

func TestSearch(t *testing.T) {
	golden := []struct {
		subPath   string
		graphPath string
		m         map[string]string
		ok        bool
	}{
		// i=0
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			m: map[string]string{
				"A": "17",
				"B": "24",
				"C": "32",
			},
			ok: true,
		},
		// i=1
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			m:         nil,
			ok:        false,
		},
		// i=2
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			m: map[string]string{
				"A": "101",
				"B": "105",
			},
			ok: true,
		},
		// i=3
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/stmt.dot",
			m: map[string]string{
				"A": "89",
				"B": "92",
				"C": "93",
			},
			ok: true,
		},

		// i=4
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			m: map[string]string{
				"A": "120",
				"B": "122",
				"C": "127",
			},
			ok: true,
		},
		// i=5
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			m: map[string]string{
				"A": "282",
				"B": "287",
				"C": "292",
				"D": "299",
			},
			ok: true,
		},
		// i=6
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			m: map[string]string{
				"A": "109",
				"B": "119",
			},
			ok: true,
		},
		// i=7
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/expr.dot",
			m: map[string]string{
				"A": "191",
				"B": "194",
				"C": "196",
			},
			ok: true,
		},

		// i=8
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			m: map[string]string{
				"A": "195",
				"B": "200",
				"C": "205",
			},
			ok: true,
		},
		// i=9
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			m: map[string]string{
				"A": "28",
				"B": "39",
				"C": "44",
				"D": "46",
			},
			ok: true,
		},
		// i=10
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			m: map[string]string{
				"A": "320",
				"B": "322",
			},
			ok: true,
		},
		// i=11
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/next.dot",
			m:         nil,
			ok:        false,
		},

		// i=12
		{
			subPath:   "../testdata/primitives/if.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			m: map[string]string{
				"A": "135",
				"B": "138",
				"C": "139",
			},
			ok: true,
		},
		// i=13
		{
			subPath:   "../testdata/primitives/if_else.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			m: map[string]string{
				"A": "442",
				"B": "448",
				"C": "451",
				"D": "453",
			},
			ok: true,
		},
		// i=14
		{
			subPath:   "../testdata/primitives/list.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			m: map[string]string{
				"A": "740",
				"B": "760",
			},
			ok: true,
		},
		// i=15
		{
			subPath:   "../testdata/primitives/pre_loop.dot",
			graphPath: "../testdata/c4_graphs/main.dot",
			m: map[string]string{
				"A": "190",
				"B": "193",
				"C": "195",
			},
			ok: true,
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
		m, ok := Search(graph, sub)
		if ok != g.ok {
			t.Errorf("i=%d: ok mismatch; expected %v, got %v", i, g.ok, ok)
			continue
		}
		if !reflect.DeepEqual(m, g.m) {
			t.Errorf("i=%d: node pair mapping mismatch; expected %v, got %v", i, g.m, m)
		}
	}
}

// sameError returns true if err is represented by the string s, and false
// otherwise. Some error messages contains "file:line" prefixes and suffixes
// from external functions, e.g.
//
//    decomp.org/x/graphs/iso.Candidates (solve.go:53): error: unable to locate entry node "foo" in graph
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
