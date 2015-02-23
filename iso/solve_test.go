package iso

import (
	"reflect"
	"testing"

	"github.com/mewfork/dot"
	"github.com/mewrev/graphs"
)

func TestCandidates(t *testing.T) {
	golden := []struct {
		graphPath string
		entry     string
		subPath   string
		want      map[string]map[string]bool
	}{
		{
			graphPath: "../testdata/primitives/if_else.dot",
			entry:     "A",
			subPath:   "../testdata/primitives/if_else.dot",
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
