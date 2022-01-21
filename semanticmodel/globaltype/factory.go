package globaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
)

type GlobalTypeModelFactory struct {
	initialGtype GlobalType
}

func CreateGlobalTypeModelFactory(
	globalTypeSexpFileName string,
) (model.ModelFactory, error) {
	gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
	if err != nil {
		return nil, err
	}
	return &GlobalTypeModelFactory{initialGtype: gtype}, nil
}

func (f GlobalTypeModelFactory) MakeModelWithLogger(logger *zap.Logger) (model.Model, error) {
	semanticModel := &globalTypeSemanticModel{
		gtype:  &f.initialGtype,
		logger: logger,
	}
	return model.MakeModelWithLogger(semanticModel, logger), nil
}
