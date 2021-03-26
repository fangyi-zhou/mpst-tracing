package pedro

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
)

type pedroSemanticModel struct {
	runtime *OcamlRuntime
}

func (p pedroSemanticModel) IsTerminated() bool {
	panic("implement me")
}

func (p pedroSemanticModel) TryReduce(action model.Action) bool {
	panic("implement me")
}

func (p pedroSemanticModel) GetEnabledActions() []model.Action {
	panic("implement me")
}

func CreatePedroSemanticModel(pedrolibFileName string, protocolFileName string, protocolName string) (model.SemanticModel, error) {
	runtime, err := LoadRuntime(pedrolibFileName)
	if err != nil {
		return nil, err
	}
	err = runtime.ImportNuscrFile(protocolFileName, protocolFileName)
	if err != nil {
		return nil, err
	}
	return pedroSemanticModel{runtime: runtime}, nil
}
