package exporter

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

type PDFExporter struct {
}

func (pdfExporter *PDFExporter) New() Exporter {
	return &PDFExporter{}
}

func (pdfExporter *PDFExporter) Generate(content string) ([]byte, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "document.html")
	if err != nil {
		return nil, err
	}
	html := "<pre>" + content + "</pre>"
	_, err = part.Write([]byte(html))
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

	var pdfBytes []byte
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}
		pdfBytes = append(pdfBytes, buf[:n]...)
	}

	return pdfBytes, nil
}

func (pdfExporter *PDFExporter) Type() string {
	return "pdf"
}
