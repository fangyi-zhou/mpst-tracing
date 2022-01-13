package mixedstateglobaltype

import (
	"github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter/semanticmodel/model"
	"strings"
)

type Choice struct {
	choicer string
	choices []MixedStateGlobalType
}

func (c Choice) PossiblePrefixes() []model.Action {
	//TODO implement me
	panic("implement me")
}

func (c Choice) ConsumePrefix(message model.Action) (MixedStateGlobalType, error) {
	//TODO implement me
	panic("implement me")
}

func (c Choice) IsDone() bool {
	//TODO implement me
	panic("implement me")
}

func (c Choice) String() string {
	//TODO implement me
	panic("implement me")
}

func (c Choice) stringWithBuilder(builder *strings.Builder) {
	//TODO implement me
	panic("implement me")
}
