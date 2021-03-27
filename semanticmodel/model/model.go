package model

import "fmt"

type Action struct {
	Src    string
	Dest   string
	Label  string
	IsSend bool
}

func (a Action) Subject() string {
	if a.IsSend {
		return a.Src
	} else {
		return a.Dest
	}
}

func (a Action) String() string {
	var action string
	if a.IsSend {
		action = "!"
	} else {
		action = "?"
	}
	return fmt.Sprintf("%s%s%s<%s>", a.Src, action, a.Dest, a.Label)
}

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
