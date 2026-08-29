package client

import (
	"net"
	"time"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/protocol"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/bet"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/logger"
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
	BatchSize  int
}

type Client struct {
	conn      net.Conn
	config    ClientConfig
	csvIter   *CsvIterator
	csvWriter *CsvWriter
	protocol  protocol.Protocol
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

	existRecord := client.csvIter.Next()

	client.protocol.SendAgencyId(client.conn, client.config.AgencyId)

	for existRecord {

		batch := 0
		var bets []bet.Bet

		for batch < client.config.BatchSize && existRecord {
			record := client.csvIter.Record()
			bet, err := bet.NewBet(client.config.AgencyId, record)
			if err != nil {
				logger.Warn("parse-csv-record", logger.Fail, "record", record)
				continue
			}
			bets = append(bets, bet)
			existRecord = client.csvIter.Next()
			batch++
		}

		logger.Info("send-bets", logger.InProgress)

		err := client.protocol.SendBets(client.conn, &bets)
		if err != nil {
			logger.Warn("send-bets", logger.Fail)
			continue
		}

		cmd, err := client.protocol.ReadCommand(client.conn)
		if err != nil {
			logger.Warn("read-command", logger.Fail)
			continue
		}
		if cmd == protocol.CMD_SERVER_ACK {
			logger.Info("received-ack", logger.Success)
		}

	}

	client.protocol.SendCommandResults(client.conn)

	cmd, err := client.protocol.ReadCommand(client.conn)
	if err != nil {
		logger.Warn("read-command", logger.Fail)
		return err
	}
	if cmd == protocol.CMD_SERVER_RESULTS {
		results, err := client.protocol.ReadResults(client.conn)
		if err != nil {
			logger.Warn("read-results", logger.Fail)
			return err
		}

		for _, bet := range results {
			logger.Info("received-result", logger.Success, "bet", bet)
			client.csvWriter.Write(bet)
		}
	}

	client.protocol.SendCommandFinished(client.conn)
	logger.Info(mainAction, logger.Success, "agency-id", client.config.AgencyId)
	//time.Sleep(ECHO_CLIENT_MESSAGE_DELAY_MS * time.Millisecond)

	return nil
}
