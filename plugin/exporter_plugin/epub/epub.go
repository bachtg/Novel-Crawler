package main

import (
	"os"

	"github.com/bmaupin/go-epub"
	"novel_crawler/internal/repository/exporter"

	"novel_crawler/constant"
	"novel_crawler/internal/model"
)

type EPUBExporter struct {
}

func (epubExporter *EPUBExporter) New() exporter.Exporter {
	return &EPUBExporter{}
}

func (epubExporter *EPUBExporter) Generate(content string) ([]byte, error) {
	e := epub.NewEpub("Collection")

	html := "<pre>" + content + "</pre>"
	_, err := e.AddSection(html, "Content", "", "")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Failed add section",
		}
	}

	err = e.Write("epubs/temp.epub")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Failed add section",
		}
	}

	contents, err := os.ReadFile("epubs/temp.epub")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Failed read file",
		}
	}

	return contents, nil
}

func (epubExporter *EPUBExporter) Type() string {
	return "epub"
}

var Exporter EPUBExporter
