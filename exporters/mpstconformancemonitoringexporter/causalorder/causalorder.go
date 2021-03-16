package causalorder

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/globaltype"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

type myMessage struct {
	globaltype.Message
}

func (m myMessage) toIndex() messageIndex {
	return messageIndex{
		label:  m.Label,
		origin: m.Origin,
		dest:   m.Dest,
	}
}

type messageIndex struct {
	label  string
	origin string
	dest   string
}

type LocalTrace []globaltype.Message

type TraceGraph struct {
	items []globaltype.Message
	graph graph.Directed
	logger *zap.Logger
}

type graphNode struct {
	graph *TraceGraph
	idx   int64
}

func (n graphNode) ID() int64 {
	return n.idx
}

func (n graphNode) DOTID() string {
	message := n.graph.items[n.idx]
	return message.String()
}

func makeNode(traceGraph *TraceGraph, id int64) graphNode {
	return graphNode{traceGraph, id}
}

func Construct(logger *zap.Logger, traces map[string]LocalTrace) TraceGraph {
	msgGraph := simple.NewDirectedGraph()
	var idx int64 = 0
	var traceGraph = TraceGraph{}
	var items = make([]globaltype.Message, 0)
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
			mIdx := myMessage{msg}.toIndex()
			if msg.Action == "send" {
				if buf, exists := recvBuffer[mIdx]; exists {
					rIdx := buf[0]
					if len(buf) == 1 {
						delete(recvBuffer, mIdx)
					} else {
						recvBuffer[mIdx] = buf[1:]
					}
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
					if len(buf) == 1 {
						delete(sendBuffer, mIdx)
					} else {
						sendBuffer[mIdx] = buf[1:]
					}
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
			idx++
		}
	}
	traceGraph.items = items
	traceGraph.graph = msgGraph
	traceGraph.logger = logger
	return traceGraph
}

func (g TraceGraph) DotFormat() (string, error) {
	dotGraph, err := dot.Marshal(g.graph, "Messages", "", "  ")
	if err != nil {
		return "", err
	} else {
		return string(dotGraph), nil
	}
}

func (g TraceGraph) CheckProtocolConformance(gty globaltype.GlobalType) error {
	tsort, err := topo.Sort(g.graph)
	if err != nil {
		return err
	}
	for _, node := range tsort {
		msg := g.items[node.ID()]
		gty, err = gty.ConsumePrefix(msg)
		if err != nil {
			return err
		}
		g.logger.Info("Consuming prefix", zap.String("msg", msg.String()), zap.String("gtype-cont", gty.String()))
	}
	if gty.IsDone() {
		return nil
	} else {
		return errors.New("global type is not done")
	}
}
