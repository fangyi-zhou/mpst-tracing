package mixedstateglobaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"go.uber.org/zap"
)

type mixedStateGlobalTypeSemanticModel struct {
	gtype  *MixedStateGlobalType
	logger *zap.Logger
}

func (g *mixedStateGlobalTypeSemanticModel) IsTerminated() bool {
	return (*g.gtype).IsDone()
}

func (g *mixedStateGlobalTypeSemanticModel) TryReduce(action model.Action) bool {
	next, err := (*g.gtype).ConsumePrefix(g, action)
	if err != nil {
		return false
	}
	g.gtype = &next
	return true
}

func (g *mixedStateGlobalTypeSemanticModel) GetEnabledActions() []model.Action {
	return (*g.gtype).PossiblePrefixes()
}

func (g *mixedStateGlobalTypeSemanticModel) SetLogger(logger *zap.Logger) {
	g.logger = logger
}

func (g *mixedStateGlobalTypeSemanticModel) Shutdown() {
	// Do nothing
}

func CreateMixedStateGlobalTypeSemanticModel(
	globalTypeSexpFileName string,
	logger *zap.Logger,
) (model.SemanticModel, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return &mixedStateGlobalTypeSemanticModel{gtype: &gtype, logger: logger}, nil
}
