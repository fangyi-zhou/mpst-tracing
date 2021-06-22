package mpstmetadatataggingprocessor

import "go.opentelemetry.io/collector/config"

type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	Roles                    map[string]metadataTag `mapstructure:"roles"`
}

type metadataTag struct {
	Name string `mapstructure:"name"`
}
