package mpstconformancecheckingprocessor

import "go.opentelemetry.io/collector/config/configmodels"

// Config defines configuration for Attributes processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
	// TODO: Add other config when needed
}
