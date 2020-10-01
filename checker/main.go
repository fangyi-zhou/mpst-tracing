package main

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenterror"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/processor/samplingprocessor/tailsamplingprocessor"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.opentelemetry.io/collector/service"
	"log"
)

func components() (component.Factories, error) {
	var errs []error

	exporters, err := component.MakeExporterFactoryMap(
		otlpexporter.NewFactory(),
	)

	if err != nil {
		errs = append(errs, err)
	}

	receivers, err := component.MakeReceiverFactoryMap(
		otlpreceiver.NewFactory(),
	)

	if err != nil {
		errs = append(errs, err)
	}

	processors, err := component.MakeProcessorFactoryMap(
		tailsamplingprocessor.NewFactory(),
		// TODO: New MPST Processor
	)

	if err != nil {
		errs = append(errs, err)
	}

	factories := component.Factories{
		Extensions: map[configmodels.Type]component.ExtensionFactory{},
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, componenterror.CombineErrors(errs)
}
func main() {
	// https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/master/cmd/otelcontribcol/main.go
	factories, err := components()
	if err != nil {
		log.Fatalf("failed to build components: %v", err)
	}

	info := component.ApplicationStartInfo{
		ExeName:  "MPSTConformance",
		LongName: "MPST Conformance Checker",
		Version:  "0.0.1-alpha",
		GitHash:  "HEAD",
	}

	params := service.Parameters{
		Factories:            factories,
		ApplicationStartInfo: info,
	}

	app, err := service.New(params)
	if err != nil {
		log.Panicf("failed to construct the application: %v", err)
		return
	}

	err = app.Run()
	if err != nil {
		log.Panicf("application run finished with error: %v", err)
		return
	}
}
