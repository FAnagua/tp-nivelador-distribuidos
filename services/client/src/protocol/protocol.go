package protocol

import (
	"encoding/binary"
	"io"
	"strconv"

	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/bet"
	"github.com/7574-sistemas-distribuidos/tp-nivelador/src/safe_socket"
)

const (
	CMD_CLIENT_AGENCY_ID uint8 = 1
	CMD_CLIENT_BET       uint8 = 2
	CMD_SERVER_ACK       uint8 = 3
	CMD_CLIENT_FINISHED  uint8 = 4
	CMD_CLIENT_RESULTS   uint8 = 5
	CMD_SERVER_RESULTS   uint8 = 6
)

type Protocol struct {
}

func (p *Protocol) serializeByte(value uint8) []byte {
	data := make([]byte, 1)
	data[0] = value

	return data
}

func (p *Protocol) serialize2Bytes(value uint16) []byte {
	data := make([]byte, 2)
	binary.BigEndian.PutUint16(data, value)

	return data
}

func (p *Protocol) serialize4Bytes(value uint32) []byte {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, value)

	return data
}

func (p *Protocol) serializeString(value string) []byte {
	data := []byte(value)
	length := uint16(len(data))
	return append(p.serialize2Bytes(length), data...)
}

func (p *Protocol) serializeBet(bet *bet.Bet) []byte {
	var data []byte

	data = append(data, p.serializeString(bet.FirstName)...)
	data = append(data, p.serializeString(bet.LastName)...)
	data = append(data, p.serialize4Bytes(uint32(bet.Document))...)
	data = append(data, p.serializeString(bet.BirthDate)...)
	data = append(data, p.serialize4Bytes(uint32(bet.BetNumber))...)

	return data
}

func (p *Protocol) readByte(socket io.Reader) (uint8, error) {

	data, err := safe_socket.RecvAll(socket, 1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

func (p *Protocol) read2Bytes(socket io.Reader) (uint16, error) {
	data, err := safe_socket.RecvAll(socket, 2)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(data), nil
}

func (p *Protocol) read4Bytes(socket io.Reader) (uint32, error) {
	data, err := safe_socket.RecvAll(socket, 4)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(data), nil
}

func (p *Protocol) readString(socket io.Reader) (string, error) {
	length, err := p.read2Bytes(socket)
	if err != nil {
		return "", err
	}

	data, err := safe_socket.RecvAll(socket, int(length))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (p *Protocol) readBet(socket io.Reader) (*bet.Bet, error) {
	firstName, err := p.readString(socket)
	if err != nil {
		return nil, err
	}

	lastName, err := p.readString(socket)
	if err != nil {
		return nil, err
	}

	document, err := p.read4Bytes(socket)
	if err != nil {
		return nil, err
	}

	birthDate, err := p.readString(socket)
	if err != nil {
		return nil, err
	}

	betNumber, err := p.read4Bytes(socket)
	if err != nil {
		return nil, err
	}

	bet := &bet.Bet{
		FirstName: firstName,
		LastName:  lastName,
		Document:  int(document),
		BirthDate: birthDate,
		BetNumber: int(betNumber),
	}
	return bet, nil
}

func (p *Protocol) ReadCommand(socket io.Reader) (uint8, error) {
	cmd, err := p.readByte(socket)
	if err != nil {
		return 0, err
	}
	return cmd, nil
}

func (p *Protocol) ReadResults(socket io.Reader) ([]*bet.Bet, error) {
	numBets, err := p.read2Bytes(socket)
	if err != nil {
		return nil, err
	}

	bets := make([]*bet.Bet, 0, numBets)
	for i := 0; i < int(numBets); i++ {
		bet, err := p.readBet(socket)
		if err != nil {
			return nil, err
		}
		bets = append(bets, bet)
	}

	return bets, nil
}

func (p *Protocol) SendAgencyId(socket io.Writer, id string) error {
	var data []byte

	value, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	cmd := p.serializeByte(CMD_CLIENT_AGENCY_ID)
	agencyId := p.serializeByte(uint8(value))
	data = append(data, cmd...)
	data = append(data, agencyId...)
	err = safe_socket.SendAll(socket, data)
	return err
}

func (p *Protocol) SendBets(socket io.Writer, bets *[]bet.Bet) error {
	var data []byte
	cmd := p.serializeByte(CMD_CLIENT_BET)
	numBets := p.serialize2Bytes(uint16(len(*bets)))
	data = append(data, cmd...)
	data = append(data, numBets...)
	for _, bet := range *bets {
		betData := p.serializeBet(&bet)
		data = append(data, betData...)
	}
	err := safe_socket.SendAll(socket, data)
	return err
}

func (p *Protocol) SendCommandFinished(socket io.Writer) error {
	cmd := p.serializeByte(CMD_CLIENT_FINISHED)
	err := safe_socket.SendAll(socket, cmd)
	return err
}

func (p *Protocol) SendCommandResults(socket io.Writer) error {
	cmd := p.serializeByte(CMD_CLIENT_RESULTS)
	err := safe_socket.SendAll(socket, cmd)
	return err
}
