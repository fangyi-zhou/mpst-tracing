package mpstmetadatataggingprocessor

import "go.opentelemetry.io/collector/config"

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
}
