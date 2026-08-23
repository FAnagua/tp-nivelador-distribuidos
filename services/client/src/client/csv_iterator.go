package client

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvIterator struct {
	file   *os.File
	reader *csv.Reader
	record []string
}

func NewCsvIterator(filePath string) (*CsvIterator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	return &CsvIterator{
		file:   file,
		reader: reader,
	}, nil
}

func (it *CsvIterator) Next() bool {
	record, err := it.reader.Read()

	if err == io.EOF {
		return false
	}

	if err != nil {
		return false
	}

	it.record = record
	return true
}

func (it *CsvIterator) Record() []string {
	return it.record
}

func (it *CsvIterator) Close() error {
	return it.file.Close()
}
