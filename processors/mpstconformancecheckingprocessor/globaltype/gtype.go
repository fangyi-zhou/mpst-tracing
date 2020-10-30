package globaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor/types"
)

type GlobalType interface {
	PossiblePrefixes() []types.Message
	ConsumePrefix(message types.Message) (GlobalType, error)
	IsDone() bool
}

func Parse(input string) (GlobalType, error) {
	return nil, errors.New("unimplemented: parse")
}
