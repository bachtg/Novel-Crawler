package repository

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Exporter interface {
	Generate(content string) ([]byte, error)
}

type ExporterManager struct {
	listExporters []*Exporter
}

type PDFExporter struct {
}

func (pdfExporter *PDFExporter) Generate(content string) ([]byte, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "document.html")
	if err != nil {
		return nil, err
	}

	_, err = part.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:3000/forms/libreoffice/convert", writer.FormDataContentType(), body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	pdfBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}
