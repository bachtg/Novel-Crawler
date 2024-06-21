package exporter

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
)

type Exporter interface {
	Generate(content string) ([]byte, error)
	Type() string
	New() Exporter
}

type ExporterManager struct {
	ExporterMapping map[string]*Exporter
}

func (exporterManager *ExporterManager) RemoveExporter(exporter string) error {
	if exporterManager.ExporterMapping[exporter] == nil {
		return &model.Err{
			Code:    constant.InternalError,
			Message: "Exporter not found",
		}
	}
	delete(exporterManager.ExporterMapping, exporter)
	return nil
}

func (exporterManager *ExporterManager) AddNewExporter(exporter ...*Exporter) error {
	if exporterManager.ExporterMapping == nil {
		exporterManager.ExporterMapping = make(map[string]*Exporter)
	}
	for _, exp := range exporter {
		exporterManager.ExporterMapping[(*exp).Type()] = exp
	}
	return nil
}
