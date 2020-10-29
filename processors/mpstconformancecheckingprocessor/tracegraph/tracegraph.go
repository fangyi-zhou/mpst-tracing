package tracegraph

import (
	"fmt"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

type Message struct {
	Label string
	Origin string
	Dest string
	Action string
}

type LocalTrace []Message

type TraceGraph struct {
	items []Message
	graph graph.Directed
}

type graphNode struct {
	graph *TraceGraph
	idx int64
}

func (n graphNode) ID() int64 {
	return n.idx
}

func (n graphNode) DOTID() string {
	message := n.graph.items[n.idx]
	return fmt.Sprintf("%s %s %s: %s", message.Origin, message.Action, message.Dest, message.Label)
}

func makeNode(traceGraph *TraceGraph, id int64) graphNode {
	return graphNode{traceGraph, id}
}

func Construct(traces map[string]LocalTrace) TraceGraph {
	msgGraph := simple.NewDirectedGraph()
	var idx int64 = 0
	var traceGraph TraceGraph
	var items = make([]Message, 0)
	for _, localTrace := range traces {
		items = append(items, localTrace[0])
		msgGraph.AddNode(makeNode(&traceGraph, idx))
		idx ++
		if len(localTrace) > 1 {
			for i := 1; i < len(localTrace); i++ {
				items = append(items, localTrace[i])
				msgGraph.AddNode(makeNode(&traceGraph, idx))
				edge := msgGraph.NewEdge(makeNode(&traceGraph, idx-1), makeNode(&traceGraph, idx))
				msgGraph.SetEdge(edge)
				idx ++
			}
		}
	}
	traceGraph.items = items
	// TODO: Add Send/Recv Edges
	// TODO: Remove debug printings
	dotGraph, err := dot.Marshal(msgGraph, "Messages", "", "  ")
	if err != nil {
		panic("unable to export to DOT")
	} else {
		fmt.Print(string(dotGraph))
	}
	return traceGraph
}