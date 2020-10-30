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

func TwoBuyer() GlobalType {
	return Send{
		origin: "A",
		dest:   "B",
		conts: map[string]GlobalType{
			"query": Recv{
				origin: "A",
				dest:   "B",
				label:  "query",
				cont:   Send{
					origin: "B",
					dest:   "A",
					conts: map[string]GlobalType{
						"quote": Recv{
							origin: "B",
							dest:   "A",
							label:  "quote",
							cont:   Send{
								origin: "B",
								dest:   "C",
								conts: map[string]GlobalType{
									"quote": Recv{
										origin: "B",
										dest:   "C",
										label:  "quote",
										cont:  Send{
											origin: "C",
											dest:   "A",
											conts: map[string]GlobalType{
												"share": Recv{
													origin: "C",
													dest:   "A",
													label:  "share",
													cont:  Send{
														origin: "A",
														dest:   "B",
														conts: map[string]GlobalType{
															"buy": Recv{
																origin: "A",
																dest:   "B",
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