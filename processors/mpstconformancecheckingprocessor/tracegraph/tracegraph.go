package tracegraph

type Message struct {
	Label string
	Origin string
	Dest string
	Action string
}

type LocalTrace []Message

type TraceGraph struct {

}

func Construct(traces map[string]LocalTrace) TraceGraph {
	return TraceGraph{}
}