package pedro

import (
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"go.uber.org/zap"
	"log"
)

type pedroSemanticModel struct {
	runtime *OcamlRuntime
	logger  *zap.Logger
}

func (p *pedroSemanticModel) IsTerminated() bool {
	return p.runtime.HasFinished()
}

func (p *pedroSemanticModel) TryReduce(action model.Action) bool {
	actionString := action.String()
	err := p.runtime.DoTransition(actionString)
	return err == nil
}

func (p *pedroSemanticModel) GetEnabledActions() []model.Action {
	transitions := p.runtime.GetEnabledTransitions()
	actions := make([]model.Action, len(transitions))
	for _, transitionString := range transitions {
		action, err := model.NewActionFromString(transitionString)
		if err != nil {
			log.Printf("internal error: unable to parse action: %s", transitionString)
			continue
		}
		actions = append(actions, action)
	}
	return actions
}

func (p *pedroSemanticModel) SetLogger(logger *zap.Logger) {
	p.logger = logger
}

func CreatePedroSemanticModel(
	pedrolibFileName string,
	protocolFileName string,
	protocolName string,
	logger *zap.Logger,
) (model.SemanticModel, error) {
	runtime, err := LoadRuntime(pedrolibFileName)
	if err != nil {
		return nil, err
	}
	logger.Info("Loaded Pedro Runtime")
	err = runtime.ImportNuscrFile(protocolFileName, protocolName)
	if err != nil {
		return nil, err
	}
	logger.Info("Imported Nuscr File", zap.String("filename", protocolName), zap.String("protocol_name", protocolName))
	return &pedroSemanticModel{runtime: runtime, logger: logger}, nil
}
