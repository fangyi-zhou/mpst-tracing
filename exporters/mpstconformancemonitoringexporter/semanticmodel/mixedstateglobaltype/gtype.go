package mixedstateglobaltype

import (
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
)

type MixedStateGlobalType interface {
	PossiblePrefixes() []model.Action
	ConsumePrefix(model *mixedStateGlobalTypeSemanticModel, message model.Action) (MixedStateGlobalType, error)
	IsDone() bool
	String() string
	ResidualActions(choicer string) [][]model.Action

	stringWithBuilder(*strings.Builder)
}
