package mpstconformancemonitoringprocessor

import "go.opentelemetry.io/collector/config"

// Config defines configuration for Attributes processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"`
	//SemanticModelType determines which semantic model is used, valid options are gtype_lts or gtype_pedro
	SemanticModelType string `mapstructure:"semantic_model_type"`
	//GlobalTypeSexpFileName is the path to a validated global protocol file, in s-expression form, used if gtype_lts model is used
	GlobalTypeSexpFileName string `mapstructure:"protocol_sexp_filename"`
	//ProtocolFileName is the path to a nuScr protocol file, used if gtype_pedro is used
	ProtocolFileName string `mapstructure:"protocol_filename"`
	//ProtocolName is the name of global protocol in file specified at ProtocolFileName
	ProtocolName string `mapstructure:"protocol_name"`
	//PedroSoFileName is the path to pedrolib.so
	PedroSoFileName string `mapstructure:"pedro_so_filename"`
}
