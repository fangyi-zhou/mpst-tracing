package model

import "go.uber.org/zap"

type SemanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetEnabledActions() []Action
}

type Model struct {
	SemanticModel
	logger *zap.Logger
	traces map[string][]Action
}

func MakeModel(semanticModel SemanticModel) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: zap.NewNop()}
}

func MakeModelWithLogger(semanticModel SemanticModel, logger *zap.Logger) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: logger}
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traces[participant] = append(m.traces[participant], traces...)
	m.logger.Info("AcceptTrace", zap.String("participant", participant), zap.Int("number", len(traces)))
	m.processTraces()
}

func (m *Model) processTraces() {
	// TODO
}
