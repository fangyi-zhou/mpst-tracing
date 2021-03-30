package model

import "go.uber.org/zap"

type SemanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetEnabledActions() []Action
}

type Model struct {
	SemanticModel
	logger  *zap.Logger
	traces  map[string][]Action
	isStuck bool
}

func MakeModel(semanticModel SemanticModel) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: zap.NewNop(), isStuck: false}
}

func MakeModelWithLogger(semanticModel SemanticModel, logger *zap.Logger) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: logger, isStuck: false}
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traces[participant] = append(m.traces[participant], traces...)
	m.logger.Info("AcceptTrace", zap.String("participant", participant), zap.Int("number", len(traces)))
	m.processTraces()
}

func (m *Model) processTraces() {
	m.logger.Info("Processing Traces")
	allHaveData := false
	for {
		allHaveData = true
		reduced := false
		for participant, trace := range m.traces {
			if len(trace) == 0 {
				allHaveData = false
				continue
			}
			action := trace[0]
			if m.TryReduce(action) {
				m.logger.Info("Action reduced successfully", zap.String("action", action.String()))
				m.traces[participant] = m.traces[participant][1:]
				reduced = true
			}
		}
		if !reduced {
			break
		}
	}
	if allHaveData && !m.IsTerminated() {
		m.logger.Error("Model is stuck")
		m.isStuck = true
	}
}

func (m *Model) IsStuck() bool {
	return m.isStuck
}
