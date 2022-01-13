package mixedstateglobaltype

import (
	"errors"
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

func Parse(input string) (MixedStateGlobalType, error) {
	return nil, errors.New("unimplemented: parse")
}

func TwoBuyer() MixedStateGlobalType {
	return Send{
		origin: "A",
		dest:   "S",
		conts: map[string]MixedStateGlobalType{
			"query": Recv{
				origin: "A",
				dest:   "S",
				label:  "query",
				cont: Send{
					origin: "S",
					dest:   "A",
					conts: map[string]MixedStateGlobalType{
						"quote": Recv{
							origin: "S",
							dest:   "A",
							label:  "quote",
							cont: Send{
								origin: "S",
								dest:   "B",
								conts: map[string]MixedStateGlobalType{
									"quote": Recv{
										origin: "S",
										dest:   "B",
										label:  "quote",
										cont: Send{
											origin: "B",
											dest:   "A",
											conts: map[string]MixedStateGlobalType{
												"share": Recv{
													origin: "B",
													dest:   "A",
													label:  "share",
													cont: Send{
														origin: "A",
														dest:   "S",
														conts: map[string]MixedStateGlobalType{
															"buy": Recv{
																origin: "A",
																dest:   "S",
																label:  "buy",
																cont:   Done{},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
