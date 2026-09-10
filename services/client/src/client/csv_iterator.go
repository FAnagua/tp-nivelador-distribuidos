package client

import (
	"bufio"
	"os"
	"strings"
)

type CsvIterator struct {
	file   *os.File
	reader *bufio.Scanner
	record []string
}

func NewCsvIterator(filePath string) (*CsvIterator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &CsvIterator{
		file:   file,
		reader: bufio.NewScanner(file),
	}, nil
}

func (it *CsvIterator) Next() bool {
	if !it.reader.Scan() {
		return false
	}

	record := it.reader.Text()

	it.record = strings.Split(record, ",")
	return true
}

func (it *CsvIterator) Record() []string {
	return it.record
}

func (it *CsvIterator) Close() error {
	return it.file.Close()
}
