package mixedstateglobaltype

import (
	model2 "github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
)

type mixedStateGlobalTypeSemanticModel struct {
	gtype  *MixedStateGlobalType
	logger *zap.Logger
	// It is important that all elements in the array must be non-empty
	residualActions [][]model2.Action
}

func (g *mixedStateGlobalTypeSemanticModel) IsTerminated() bool {
	return (*g.gtype).IsDone()
}

func (g *mixedStateGlobalTypeSemanticModel) TryReduce(action model2.Action) bool {
	// First see if any residual actions can reduce
	// They should be disjoint from main actions (hopefully...)
	for idx, actions := range g.residualActions {
		if actions[0] == action {
			if len(actions) == 1 {
				// Clean up by moving the last in the array to current position
				g.residualActions[idx] = g.residualActions[len(g.residualActions)-1]
				g.residualActions = g.residualActions[:len(g.residualActions)-1]
			} else {
				g.residualActions[idx] = actions[1:]
			}
			g.logger.Info("Reduced residual action", zap.String("action", action.String()))
			return true
		}
	}
	g.logger.Info("Global Type", zap.String("global", (*g.gtype).String()))
	next, err := (*g.gtype).ConsumePrefix(g, action)
	if err != nil {
		return false
	}
	g.gtype = &next
	return true
}

func (g *mixedStateGlobalTypeSemanticModel) GetEnabledActions() []model2.Action {
	return (*g.gtype).PossiblePrefixes()
}

func (g *mixedStateGlobalTypeSemanticModel) SetLogger(logger *zap.Logger) {
	g.logger = logger
}

func (g *mixedStateGlobalTypeSemanticModel) Shutdown() {
	// Do nothing
}

func (g *mixedStateGlobalTypeSemanticModel) AddResidualActions(residuals [][]model2.Action) {
	g.residualActions = append(g.residualActions, residuals...)
}

func CreateMixedStateGlobalTypeSemanticModel(
	globalTypeSexpFileName string,
	logger *zap.Logger,
) (model2.SemanticModel, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return &mixedStateGlobalTypeSemanticModel{gtype: &gtype, logger: logger}, nil
}
