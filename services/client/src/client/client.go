package client

import (
	"net"
	"strings"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const CONNECTION_ATTEMPTS_MAX = 3
const CONNECTION_ATTEMPS_DELAY_MS = 200

const ECHO_CLIENT_BUFFER_SIZE = 512
const ECHO_CLIENT_MESSAGE_AMOUNT = 3
const ECHO_CLIENT_MESSAGE_DELAY_MS = 1000

type ClientConfig struct {
	ServerHost string
	ServerPort string
	AgencyId   string
	InputFile  string
	OutputFile string
}

type Client struct {
	conn      net.Conn
	config    ClientConfig
	csvIter   *CsvIterator
	csvWriter *CsvWriter
}

func NewClient(config ClientConfig) (*Client, error) {
	conn, err := connectToServer(config.ServerHost, config.ServerPort)
	if err != nil {
		logger.Warn("connect-to-server", logger.Fail)
		return nil, err
	}

	csvIter, err := NewCsvIterator(config.InputFile)
	if err != nil {
		logger.Warn("open-csv-file", logger.Fail)
		return nil, err
	}

	csvWriter, err := NewCsvWriter(config.OutputFile)
	if err != nil {
		logger.Warn("open-output-file", logger.Fail)
		return nil, err
	}

	client := &Client{conn: conn, config: config, csvIter: csvIter, csvWriter: csvWriter}
	return client, nil
}

func connectToServer(host, port string) (net.Conn, error) {
	const action = "connect-to-server"
	var err error
	var conn net.Conn

	logger.Info(action, logger.InProgress)
	for i := range CONNECTION_ATTEMPTS_MAX {
		conn, err = net.Dial("tcp", host+":"+port)
		if err != nil {
			logger.Warn(action, logger.Fail, "attempt", i)
			time.Sleep(CONNECTION_ATTEMPS_DELAY_MS * time.Millisecond)
			continue
		}

		logger.Info(action, logger.Success)
		break
	}

	return conn, err
}

func (client *Client) Run() error {
	const mainAction = "test-echo-server"
	defer client.conn.Close()
	defer client.csvIter.Close()
	defer client.csvWriter.Close()

	var record []string
	var messageId string

	for client.csvIter.Next() {
		record = client.csvIter.Record()
		messageId = record[0]

		messageArgs := []any{"agency-id", client.config.AgencyId, "message-id", messageId}
		logger.Info(mainAction, logger.InProgress, messageArgs...)

		logger.Info("read-csv-record", logger.InProgress, "record", record)

		clientMessage := record[0] + "," + record[1] + "," + record[2] + "," + record[3] + "," + record[4]

		if err := safe_socket.SendAll(client.conn, []byte(clientMessage)); err != nil {
			logger.Error("send-message", logger.Fail, messageArgs...)
			return err
		}

		responseBuffer, err := safe_socket.RecvAll(client.conn, ECHO_CLIENT_BUFFER_SIZE)
		if err != nil {
			logger.Error("recv-response", logger.Fail, messageArgs...)
			return err
		}

		if string(responseBuffer) == clientMessage {
			logger.Error("check-response", logger.Fail, messageArgs...)
			return err
		}

		responseBuffer = responseBuffer[:len(clientMessage)]

		if err := client.csvWriter.Write(strings.Split(string(responseBuffer), ",")); err != nil {
			logger.Error("write-csv-record", logger.Fail, messageArgs...)
			return err
		}

		time.Sleep(ECHO_CLIENT_MESSAGE_DELAY_MS * time.Millisecond)
	}
	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)

	return nil
}
