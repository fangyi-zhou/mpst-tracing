package model

import (
	"go.uber.org/zap"
	"sync"
)

type modelState int

//go:generate stringer -type=modelState
const (
	NORMAL modelState = iota
	STUCK
	TERMINATED
)

type SemanticModel interface {
	IsTerminated() bool
	TryReduce(action Action) bool
	GetEnabledActions() []Action
	SetLogger(logger *zap.Logger)
	Shutdown()
}

type Model struct {
	SemanticModel
	logger    *zap.Logger
	traces    map[string][]Action
	state     modelState
	traceLock *sync.Mutex
}

func MakeModel(semanticModel SemanticModel) Model {
	logger := zap.NewNop()
	return MakeModelWithLogger(semanticModel, logger)
}

func MakeModelWithLogger(semanticModel SemanticModel, logger *zap.Logger) Model {
	return Model{
		SemanticModel: semanticModel,
		traces:        make(map[string][]Action),
		logger:        logger,
		state:         NORMAL,
		traceLock:     &sync.Mutex{},
	}
}

func (m *Model) AcceptTrace(participant string, traces []Action) {
	m.traceLock.Lock()
	defer m.traceLock.Unlock()
	m.traces[participant] = append(m.traces[participant], traces...)
	m.logger.Info(
		"AcceptTrace",
		zap.String("participant", participant),
		zap.Int("number", len(traces)),
	)
	m.processTraces()
}

func (m *Model) processTraces() {
	m.logger.Info("Processing Traces", zap.String("model-state", m.state.String()))
	if m.state != NORMAL {
		// No need to process if model is stuck or terminated
		return
	}
	for {
		reduced := false
		for participant, trace := range m.traces {
			if len(trace) == 0 {
				m.logger.Info(
					"Trace queue is empty for participant",
					zap.String("participant", participant),
				)
				continue
			}
			action := trace[0]
			if m.TryReduce(action) {
				m.logger.Info(
					"Action reduced successfully",
					zap.String("action", action.String()),
					zap.String("participant", participant),
				)
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
	m.updateStatus()
}

func (m *Model) updateStatus() {
	// assumed state is NORMAL when calling this function
	if m.SemanticModel.IsTerminated() {
		m.state = TERMINATED
		m.logger.Info("Model terminated")
		return
	}
	if m.isStuck() {
		m.state = STUCK
		m.logger.Error("Model is stuck!")
		return
	}
	// model is still NORMAL
}

// Determines whether a model is stuck by checking the enabled actions
func (m *Model) isStuck() bool {
	enabledActions := m.SemanticModel.GetEnabledActions()
	// enabledActionStrings := make([]string, 0)
	// for _, action := range enabledActions {
	// 	enabledActionStrings = append(enabledActionStrings, action.String())
	// }
	// m.logger.Info("Enabled actions", zap.Strings("actions", enabledActionStrings))
	actionsByParticipant := make(map[string][]Action)
	// Group enabled actions by their subjects
	for _, action := range enabledActions {
		subject := action.Subject()
		actionsByParticipant[subject] = append(actionsByParticipant[subject], action)
	}
	for subject, actions := range actionsByParticipant {
		if len(m.traces[subject]) == 0 {
			// The model is not necessarily stuck because the message may be yet to arrive
			return false
		}
		firstAction := m.traces[subject][0]
		for _, action := range actions {
			if action == firstAction {
				return false
			}
		}
	}
	return true
}

func (m *Model) IsStuck() bool {
	return m.state == STUCK
}
