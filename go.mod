module github.com/fangyi-zhou/mpst-tracing

go 1.15

require (
	github.com/fangyi-zhou/mpst-tracing/twobuyer v0.0.0
	github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor v0.0.0
)

replace (
	github.com/fangyi-zhou/mpst-tracing/twobuyer => ./twobuyer
	github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor => ./processors/mpstconformancecheckingprocessor
)
