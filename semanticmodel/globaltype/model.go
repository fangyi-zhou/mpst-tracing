package globaltype

import "github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"

type globalTypeSemanticModel struct {
	gtype GlobalType
}

func (g globalTypeSemanticModel) IsTerminated() bool {
	return g.gtype.IsDone()
}

func (g globalTypeSemanticModel) TryReduce(action model.Action) bool {
	panic("implement me")
}

func (g globalTypeSemanticModel) GetEnabledActions() []model.Action {
	panic("implement me")
}

func CreateGlobalTypeSemanticModel(globalTypeSexpFileName string) (model.SemanticModel, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return globalTypeSemanticModel{gtype: gtype}, nil
}
