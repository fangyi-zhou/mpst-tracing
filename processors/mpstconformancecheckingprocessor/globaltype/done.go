package globaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor/types"
)

type Done struct{}

func (Done) PossiblePrefixes() []types.Message {
	return nil
}

func (Done) ConsumePrefix(message types.Message) (GlobalType, error) {
	return nil, errors.New("done cannot consume prefix")
}

func (Done) isDone() bool {
	return true
}
