package iso

import (
	"reflect"
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
	}{
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
		},
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
		},
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
		},
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
		},
	}

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
		if !reflect.DeepEqual(eq.c, g.want) {
			t.Errorf("i=%d: candidate map mismatch; expected %v, got %v", i, g.want, eq.c)
		}
	}
}
