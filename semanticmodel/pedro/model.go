package pedro

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"log"
)

type pedroSemanticModel struct {
	runtime *OcamlRuntime
}

func (p pedroSemanticModel) IsTerminated() bool {
	return p.runtime.HasFinished()
}

func (p pedroSemanticModel) TryReduce(action model.Action) bool {
	actionString := action.String()
	err := p.runtime.DoTransition(actionString)
	return err != nil
}

func (p pedroSemanticModel) GetEnabledActions() []model.Action {
	transitions := p.runtime.GetEnabledTransitions()
	actions := make([]model.Action, len(transitions))
	for _, transitionString := range transitions {
		action, err := model.NewActionFromString(transitionString)
		if err != nil {
			log.Panicf("internal error: unable to parse action: %s", transitionString)
		}
		actions = append(actions, action)
	}
	return actions
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
