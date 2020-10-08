module github.com/fangyi-zhou/mpst-tracing

go 1.15

require (
	github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor v0.0.0
	go.opentelemetry.io/collector v0.11.0
	go.opentelemetry.io/otel v0.12.0
	go.opentelemetry.io/otel/exporters/otlp v0.12.0
	go.opentelemetry.io/otel/exporters/stdout v0.12.0
	go.opentelemetry.io/otel/exporters/trace/jaeger v0.12.0
	go.opentelemetry.io/otel/sdk v0.12.0
	golang.org/x/net v0.0.0-20200930145003-4acb6c075d10 // indirect
	golang.org/x/sync v0.0.0-20200930132711-30421366ff76 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	google.golang.org/genproto v0.0.0-20201001141541-efaab9d3c4f7 // indirect
)

replace github.com/fangyi-zhou/mpst-tracing/processors/mpstconformancecheckingprocessor => ./processors/mpstconformancecheckingprocessor
