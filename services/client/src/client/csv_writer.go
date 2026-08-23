package client

import (
	"encoding/csv"
	"os"
)

type CsvWriter struct {
	file   *os.File
	writer *csv.Writer
}

func NewCsvWriter(filePath string) (*CsvWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)

	return &CsvWriter{
		file:   file,
		writer: writer,
	}, nil
}

func (w *CsvWriter) Write(record []string) error {
	err := w.writer.Write(record)
	if err != nil {
		return err
	}
	w.writer.Flush()
	return w.writer.Error()
}

func (w *CsvWriter) Close() error {
	w.writer.Flush()
	err := w.writer.Error()
	if err != nil {
		return err
	}
	return w.file.Close()
}
