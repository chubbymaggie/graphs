## WIP

This project is a *work in progress*. The implementation is *incomplete* and subject to change. The documentation may be inaccurate.

# graphs

[![GoDoc](https://godoc.org/github.com/mewrev/graphs?status.svg)](https://godoc.org/github.com/mewrev/graphs)

The graphs project implements subgraph isomorphism search algorithms.

## cmd/iso

`iso` is a tool which locates subgraph isomorphisms in graphs.

### Installation

```shell
go get github.com/mewrev/graphs/cmd/iso
```

### Usage

    Usage: iso [OPTION]... SUB.dot GRAPH.dot

    Flags:
      -all=true: Locate all isomorphisms of SUB in GRAPH.

### Examples

1) Locate all isomorphisms of the subgraph [if.dot](testdata/if.dot) in the graph [stmt.dot](testdata/c4_graphs/stmt.dot).

```bash
iso if.dot c4_graphs/stmt.dot
// Output:
// Isomorphism found at node "17":
//    "A"="17"
//    "B"="24"
//    "C"="32"
// Isomorphism found at node "71":
//    "A"="71"
//    "B"="74"
//    "C"="75"
```

SUB:
* [if.dot](testdata/if.dot):

![if.dot subgraph](https://raw.githubusercontent.com/mewrev/graphs/master/testdata/if.png)

GRAPH:
* [stmt.dot](testdata/c4_graphs/stmt.dot):

![stmt.dot graph](https://raw.githubusercontent.com/mewrev/graphs/master/testdata/c4_graphs/stmt.png)

## Public domain

The source code and any original content of this repository is hereby released into the [public domain].

[public domain]: https://creativecommons.org/publicdomain/zero/1.0/
