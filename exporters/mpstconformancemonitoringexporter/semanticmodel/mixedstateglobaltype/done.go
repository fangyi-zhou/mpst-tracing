package mixedstateglobaltype

import (
	"errors"
	"strings"

	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
)

type Done struct{}

func (Done) PossiblePrefixes() []model.Action {
	return make([]model.Action, 0)
}

func (Done) ConsumePrefix(_ model.Action) (MixedStateGlobalType, error) {
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

func NewDone() Done {
	return Done{}
}
