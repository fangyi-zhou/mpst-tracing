package globaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"strings"
)

type GlobalType interface {
	PossiblePrefixes() []model.Action
	ConsumePrefix(message model.Action) (GlobalType, error)
	IsDone() bool
	String() string

	stringWithBuilder(*strings.Builder)
}
