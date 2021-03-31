package model

import (
	"go.uber.org/zap"
	"sync"
)

type SemanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetEnabledActions() []Action
}

type Model struct {
	SemanticModel
	logger    *zap.Logger
	traces    map[string][]Action
	isStuck   bool
	traceLock *sync.Mutex
}

func MakeModel(semanticModel SemanticModel) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: zap.NewNop(), isStuck: false, traceLock: &sync.Mutex{}}
}

func MakeModelWithLogger(semanticModel SemanticModel, logger *zap.Logger) Model {
	return Model{SemanticModel: semanticModel, traces: make(map[string][]Action), logger: logger, isStuck: false, traceLock: &sync.Mutex{}}
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traceLock.Lock()
	defer m.traceLock.Unlock()
	m.traces[participant] = append(m.traces[participant], traces...)
	m.logger.Info("AcceptTrace", zap.String("participant", participant), zap.Int("number", len(traces)))
	m.processTraces()
}

func (m *Model) processTraces() {
	m.logger.Info("Processing Traces")
	for {
		reduced := false
		for participant, trace := range m.traces {
			if len(trace) == 0 {
				m.logger.Info("Trace queue is empty for participant", zap.String("participant", participant))
				continue
			}
			action := trace[0]
			if m.TryReduce(action) {
				m.logger.Info("Action reduced successfully", zap.String("action", action.String()), zap.String("participant", participant))
				m.traces[participant] = m.traces[participant][1:]
				reduced = true
			} else {
				m.logger.Info("Cannot reduce action for participant", zap.String("action", action.String()), zap.String("participant", participant))
			}
		}
		if !reduced {
			break
		}
	}
	m.checkStuck()
	if m.isStuck {
		m.logger.Error("Model is stuck")
	}
}

func (m *Model) checkStuck() {
	// TODO: check whether model is stuck
}

func (m *Model) IsStuck() bool {
	return m.isStuck
}
