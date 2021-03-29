package model

type SemanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetEnabledActions() []Action
}

type Model struct {
	SemanticModel
	traces map[string][]Action
}

func MakeModel(semanticModel SemanticModel) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action)}
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traces[participant] = append(m.traces[participant], traces...)
	m.processTraces()
}

func (m *Model) processTraces() {
	// TODO
}
