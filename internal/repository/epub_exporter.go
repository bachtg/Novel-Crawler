package repository

import (
	"github.com/bmaupin/go-epub"
	"io/ioutil"
	"novel_crawler/internal/model"
	"novel_crawler/constant"
)

type EPUBExporter struct {
}

func NewEpubExporter() Exporter {
	return &EPUBExporter{}
}

func (epubExporter *EPUBExporter) Generate(content string) ([]byte, error) {
	e := epub.NewEpub("Collection")

	_, err := e.AddSection(content, "Content", "", "")
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

	contents, err := ioutil.ReadFile("epubs/temp.epub")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Failed read file",
		}
	}

	return contents, nil
}
