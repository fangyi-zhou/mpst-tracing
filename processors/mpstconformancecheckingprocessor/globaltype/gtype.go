package globaltype

import (
	"errors"
	"github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor/types"
	"strings"
)

type GlobalType interface {
	PossiblePrefixes() []types.Message
	ConsumePrefix(message types.Message) (GlobalType, error)
	IsDone() bool
	String() string

	stringWithBuilder(*strings.Builder)
}

func Parse(input string) (GlobalType, error) {
	return nil, errors.New("unimplemented: parse")
}

func TwoBuyer() GlobalType {
	return Send{
		origin: "A",
		dest:   "S",
		conts: map[string]GlobalType{
			"query": Recv{
				origin: "A",
				dest:   "S",
				label:  "query",
				cont: Send{
					origin: "S",
					dest:   "A",
					conts: map[string]GlobalType{
						"quote": Recv{
							origin: "S",
							dest:   "A",
							label:  "quote",
							cont: Send{
								origin: "S",
								dest:   "B",
								conts: map[string]GlobalType{
									"quote": Recv{
										origin: "S",
										dest:   "B",
										label:  "quote",
										cont: Send{
											origin: "B",
											dest:   "A",
											conts: map[string]GlobalType{
												"share": Recv{
													origin: "B",
													dest:   "A",
													label:  "share",
													cont: Send{
														origin: "A",
														dest:   "S",
														conts: map[string]GlobalType{
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
