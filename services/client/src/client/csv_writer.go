package client

import (
	"bufio"
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/bet"
)

type CsvWriter struct {
	file   *os.File
	writer *bufio.Writer
}

func NewCsvWriter(filePath string) (*CsvWriter, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	return &CsvWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (w *CsvWriter) Write(bet *bet.Bet) error {
	record := bet.FirstName + "," +
		bet.LastName + "," +
		strconv.Itoa(bet.Document) + "," +
		bet.BirthDate + "," +
		strconv.Itoa(bet.BetNumber) + "\n"

	_, err := w.writer.WriteString(record)

	return err
}

func (w *CsvWriter) Close() error {
	err := w.writer.Flush()
	if err != nil {
		return err
	}
	return w.file.Close()
}
