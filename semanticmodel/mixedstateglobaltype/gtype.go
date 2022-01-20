package mixedstateglobaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"strings"
)

type MixedStateGlobalType interface {
	PossiblePrefixes() []model.Action
	ConsumePrefix(model *mixedStateGlobalTypeSemanticModel, message model.Action) (MixedStateGlobalType, error)
	IsDone() bool
	String() string
	ResidualActions(choicer string) [][]model.Action

	stringWithBuilder(*strings.Builder)
}
