package mpstconformancecheckingprocessor

import "go.opentelemetry.io/collector/config/configmodels"

// Config defines configuration for Attributes processor.
type Config struct {
	configmodels.ProcessorSettings `mapstructure:",squash"`
	// Protocol is a global type in S-expression format
	Protocol string `mapstructure:"protocol"`
}
