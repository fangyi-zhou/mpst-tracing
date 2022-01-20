package mixedstateglobaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"strings"
)

type Done struct{}

func (Done) PossiblePrefixes() []model.Action {
	return make([]model.Action, 0)
}

func (Done) ConsumePrefix(_ *mixedStateGlobalTypeSemanticModel, _ model.Action) (MixedStateGlobalType, error) {
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

func (Done) ResidualActions(_ string) [][]model.Action {
	return [][]model.Action{{}}
}

func NewDone() Done {
	return Done{}
}
