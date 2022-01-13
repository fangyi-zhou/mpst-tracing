package mixedstateglobaltype

import (
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
)

type MixedStateGlobalType interface {
	PossiblePrefixes() []model.Action
	ConsumePrefix(message model.Action) (MixedStateGlobalType, error)
	IsDone() bool
	String() string

	stringWithBuilder(*strings.Builder)
}
