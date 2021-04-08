package globaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"go.uber.org/zap"
)

type globalTypeSemanticModel struct {
	gtype  *GlobalType
	logger *zap.Logger
}

func (g *globalTypeSemanticModel) IsTerminated() bool {
	return (*g.gtype).IsDone()
}

func (g *globalTypeSemanticModel) TryReduce(action model.Action) bool {
	next, err := (*g.gtype).ConsumePrefix(action)
	if err != nil {
		return false
	}
	g.gtype = &next
	return true
}

func (g *globalTypeSemanticModel) GetEnabledActions() []model.Action {
	return (*g.gtype).PossiblePrefixes()
}

func (g *globalTypeSemanticModel) SetLogger(logger *zap.Logger) {
	g.logger = logger
}

func CreateGlobalTypeSemanticModel(globalTypeSexpFileName string, logger *zap.Logger) (model.SemanticModel, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return &globalTypeSemanticModel{gtype: &gtype, logger: logger}, nil
}
