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

func (m Message) toIndex() messageIndex {
	return messageIndex{
		label: m.Label,
		origin: m.Origin,
		dest: m.Dest,
	}
}

type messageIndex struct {
	label string
	origin string
	dest string
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
	return fmt.Sprintf("%s_%s_%s:%s", message.Origin, message.Action, message.Dest, message.Label)
}

func makeNode(traceGraph *TraceGraph, id int64) graphNode {
	return graphNode{traceGraph, id}
}

func Construct(traces map[string]LocalTrace) TraceGraph {
	msgGraph := simple.NewDirectedGraph()
	var idx int64 = 0
	var traceGraph TraceGraph
	var items = make([]Message, 0)
	var sendBuffer = map[messageIndex][]int64{}
	var recvBuffer = map[messageIndex][]int64{}
	for _, localTrace := range traces {
		for i, msg := range localTrace {
			items = append(items, msg)
			msgGraph.AddNode(makeNode(&traceGraph, idx))
			if i > 0 {
				// Local Sequencing
				edge := msgGraph.NewEdge(makeNode(&traceGraph, idx-1), makeNode(&traceGraph, idx))
				msgGraph.SetEdge(edge)
			}
			mIdx := msg.toIndex()
			if msg.Action == "send" {
				if buf, exists := recvBuffer[mIdx]; exists {
					rIdx := buf[0]
					recvBuffer[mIdx] = buf[1:]
					edge := msgGraph.NewEdge(makeNode(&traceGraph, idx), makeNode(&traceGraph, rIdx))
					msgGraph.SetEdge(edge)
				} else {
					if buf, exists := sendBuffer[mIdx]; exists {
						sendBuffer[mIdx] = append(buf, idx)
					} else {
						sendBuffer[mIdx] = []int64{idx}
					}
				}
			} else if msg.Action == "recv" {
				if buf, exists := sendBuffer[mIdx]; exists {
					sIdx := buf[0]
					sendBuffer[mIdx] = buf[1:]
					edge := msgGraph.NewEdge(makeNode(&traceGraph, sIdx), makeNode(&traceGraph, idx))
					msgGraph.SetEdge(edge)
				} else {
					if buf, exists := recvBuffer[mIdx]; exists {
						recvBuffer[mIdx] = append(buf, idx)
					} else {
						recvBuffer[mIdx] = []int64{idx}
					}
				}
			} else {
				panic("unknown action " + msg.Action)
			}
			idx ++
		}
	}
	traceGraph.items = items
	// TODO: Remove debug printings
	dotGraph, err := dot.Marshal(msgGraph, "Messages", "", "  ")
	if err != nil {
		panic("unable to export to DOT")
	} else {
		fmt.Print(string(dotGraph))
	}
	return traceGraph
}