package mixedstateglobaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
)

type MixedStateGlobalTypeModelFactory struct {
	initialGtype MixedStateGlobalType
}

func CreateMixedStateGlobalTypeModelFactory(
	globalTypeSexpFileName string,
) (model.ModelFactory, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return &MixedStateGlobalTypeModelFactory{initialGtype: gtype}, nil
}

func (f MixedStateGlobalTypeModelFactory) MakeModelWithLogger(logger *zap.Logger) (model.Model, error) {
	semanticModel := &mixedStateGlobalTypeSemanticModel{
		gtype:  &f.initialGtype,
		logger: logger,
	}
	return model.MakeModelWithLogger(semanticModel, logger), nil
}
