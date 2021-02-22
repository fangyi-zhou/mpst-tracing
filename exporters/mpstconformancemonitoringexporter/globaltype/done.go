package globaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/types"
	"strings"
)

type Done struct{}

func (Done) PossiblePrefixes() []types.Message {
	return nil
}

func (Done) ConsumePrefix(message types.Message) (GlobalType, error) {
	return nil, errors.New("end cannot consume prefix")
}

func (Done) IsDone() bool {
	return true
}

func (Done) String() string {
	return "end"
}

func (Done) stringWithBuilder(b *strings.Builder) {
	b.WriteString("end")
}
