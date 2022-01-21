package pedro

import (
	"github.com/fangyi-zhou/mpst-tracing/semanticmodel/model"
	"go.uber.org/zap"
)

type PedroModelFactory struct {
	libFileName      string
	protocolFileName string
	protocolName     string
}

func CreatePedroModelFactory(
	pedrolibFileName string,
	protocolFileName string,
	protocolName string,
) (model.ModelFactory, error) {
	return &PedroModelFactory{
		libFileName:      pedrolibFileName,
		protocolFileName: protocolFileName,
		protocolName:     protocolName,
	}, nil
}

func (p PedroModelFactory) MakeModelWithLogger(logger *zap.Logger) (model.Model, error) {
	runtime, err := LoadRuntime(p.libFileName)
	if err != nil {
		return model.Model{}, err
	}
	logger.Info("Loaded Pedro Runtime")
	err = runtime.ImportNuscrFile(p.protocolFileName, p.protocolName)
	if err != nil {
		return model.Model{}, err
	}
	logger.Info(
		"Imported Nuscr File",
		zap.String("filename", p.protocolName),
		zap.String("protocol_name", p.protocolName),
	)
	semanticModel := &pedroSemanticModel{runtime: runtime, logger: logger}
	return model.MakeModelWithLogger(semanticModel, logger), nil
}
