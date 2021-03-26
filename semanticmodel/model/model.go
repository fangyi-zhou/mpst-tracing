package model

type Action struct {
	Src   string
	Dest  string
	Label string
}

type semanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetAllActions() []string
}

type Model struct {
	semanticModel
	traces map[string][]Action
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traces[participant] = append(m.traces[participant], traces...)
	m.processTraces()
}

func (m *Model) processTraces() {
	// TODO
}
