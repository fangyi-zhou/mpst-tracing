module github.com/fangyi-zhou/mpst-tracing/exporters/mpstconformancemonitoringexporter

go 1.15

require (
	github.com/fangyi-zhou/mpst-tracing/semanticmodel v0.0.0-20210330132437-4b327c0d6af1
	github.com/golang/protobuf v1.5.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.22.0
	go.uber.org/zap v1.16.0
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
)

replace github.com/fangyi-zhou/mpst-tracing/semanticmodel => ./../../semanticmodel
