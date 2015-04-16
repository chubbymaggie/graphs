// Package primitive defines the types used to represent high-level control flow
// primitives.
package primitive

// A Primitive represents a high-level control flow primitive (e.g. 2-way
// conditional, pre-test loop) as a mapping from subgraph (graph representation
// of a control flow primitive) node names to control flow graph node names.
type Primitive struct {
	// Primitive name; e.g. "if", "pre_loop", ...
	Prim string `json:"prim"`
	// Node name of the primitive; e.g. "list0".
	Node string `json:"node"`
	// Node mapping; e.g. {"A": 1, "B": 2, "C": 3}
	Nodes map[string]string `json:"nodes"`
}
