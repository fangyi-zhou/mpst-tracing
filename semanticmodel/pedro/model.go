package pedro

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
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
	// p.logger.Info("Raw enabled actions", zap.Strings("raw_actions", transitions))
	actions := make([]model.Action, 0)
	for _, transitionString := range transitions {
		if len(transitionString) > 0 && transitionString[0] == '_' {
			// Silent Transitions begin with underscore, skipping
			continue
		}
		action, err := model.NewActionFromString(transitionString)
		if err != nil {
			p.logger.Info(
				"skipping unrecognised action string",
				zap.String("raw", transitionString),
			)
			continue
		}
		actions = append(actions, action)
	}
	return actions
}

func (p *pedroSemanticModel) SetLogger(logger *zap.Logger) {
	p.logger = logger
}

func (p *pedroSemanticModel) Shutdown() {
	p.runtime.Close()
}
