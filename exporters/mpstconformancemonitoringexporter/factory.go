package mpstconformancemonitoringexporter

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr configmodels.Type = "mpstconformancemonitoring"
)

func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTraceExporter))
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			NameVal: "MPSTConformanceMonitoring",
			TypeVal: typeStr,
		},
	}
}

func createTraceExporter(
	ctx context.Context,
	params component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.TracesExporter, error) {
	return newMpstConformanceExporter(params.Logger, cfg.(*Config))
}
