package globaltype

import (
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
)

type GlobalType interface {
	PossiblePrefixes() []model.Action
	ConsumePrefix(message model.Action) (GlobalType, error)
	IsDone() bool
	String() string

	stringWithBuilder(*strings.Builder)
}
