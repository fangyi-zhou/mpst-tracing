package globaltype

import (
	"fmt"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
)

type GlobalTypeModelFactory struct {
	initialGtype GlobalType
}

func CreateGlobalTypeModelFactory(
	globalTypeSexpFileName string,
	globalTypeProtobufFileName string,
) (model.ModelFactory, error) {
	if globalTypeSexpFileName != "" {
		gtype, err := LoadFromSexpFile(globalTypeSexpFileName)
		if err != nil {
			return nil, err
		}
		return &GlobalTypeModelFactory{initialGtype: gtype}, nil
	} else if globalTypeProtobufFileName != "" {
		gtype, err := LoadFromProtobuf(globalTypeProtobufFileName)
		if err != nil {
			return nil, err
		}
		return &GlobalTypeModelFactory{initialGtype: gtype}, nil
	}
	return nil, fmt.Errorf(
		"must provide a global protocol via either a s-expression or a protobuf file",
	)
}

func (f GlobalTypeModelFactory) MakeModelWithLogger(logger *zap.Logger) (model.Model, error) {
	semanticModel := &globalTypeSemanticModel{
		gtype:  &f.initialGtype,
		logger: logger,
	}
	return model.MakeModelWithLogger(semanticModel, logger), nil
}
