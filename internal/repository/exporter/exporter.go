package exporter

import (
	"novel_crawler/internal/model"
	"novel_crawler/constant"
)

type Exporter interface {
	Generate(content string) ([]byte, error)
	Type() string
}

type ExporterManager struct {
	ExporterMapping map[string]*Exporter
}


func (exporterManager *ExporterManager) RemoveExporter(exporter string) error {
	if exporterManager.ExporterMapping[exporter] == nil {
        return &model.Err{
			Code: constant.InternalError,
            Message: "Exporter not found",
		}
    }
	delete(exporterManager.ExporterMapping, exporter)
	return nil
}