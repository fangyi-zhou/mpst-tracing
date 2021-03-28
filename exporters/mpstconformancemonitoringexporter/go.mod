module github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter

go 1.15

require (
	github.com/fangyi-zhou/mpst-tracing/semanticmodel v0.0.0-20210326173509-4b7cf9b0bb4f
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.20.0
	go.uber.org/zap v1.16.0
)

replace github.com/fangyi-zhou/mpst-tracing/semanticmodel => ./../../semanticmodel
