package iso

import (
	"log"

	"github.com/mewfork/dot"
	"github.com/mewrev/graphs"
)

// Isomorphism returns a mapping from sub node name to graph node name if there
// exists an isomorphism of sub in graph which starts at the entry node. The
// boolean value is true if such a mapping could be located, and false
// otherwise.
func Isomorphism(graph *dot.Graph, entry string, sub *graphs.SubGraph) (m map[string]string, ok bool) {
	eq, err := Candidates(graph, entry, sub)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	err = eq.Solve(graph, sub)
	if err != nil {
		log.Println(err)
		return nil, false
	}
	return eq.m, true
}
