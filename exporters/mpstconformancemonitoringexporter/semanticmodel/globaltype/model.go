package globaltype

import "github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"

type globalTypeSemanticModel struct {
	gtype *GlobalType
}

func (g globalTypeSemanticModel) IsTerminated() bool {
	return (*g.gtype).IsDone()
}

func (g globalTypeSemanticModel) TryReduce(action model.Action) bool {
	next, err := (*g.gtype).ConsumePrefix(action)
	if err != nil {
		return false
	}
	g.gtype = &next
	return true
}

func (g globalTypeSemanticModel) GetEnabledActions() []model.Action {
	return (*g.gtype).PossiblePrefixes()
}

func CreateGlobalTypeSemanticModel(globalTypeSexpFileName string) (model.SemanticModel, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return globalTypeSemanticModel{gtype: &gtype}, nil
}
