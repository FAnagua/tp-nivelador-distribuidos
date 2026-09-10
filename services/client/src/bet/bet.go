package bet

import (
	"strconv"
)

type Bet struct {
	AgencyId  int
	FirstName string
	LastName  string
	Document  int
	BirthDate string
	BetNumber int
}

func NewBet(agencyId string, record []string) (Bet, error) {
	agencyIdInt, err := strconv.Atoi(agencyId)
	if err != nil {
		return Bet{}, err
	}
	betNumber, err := strconv.Atoi(record[4])
	if err != nil {
		return Bet{}, err
	}
	document, err := strconv.Atoi(record[2])
	if err != nil {
		return Bet{}, err
	}

	return Bet{
		AgencyId:  agencyIdInt,
		FirstName: record[0],
		LastName:  record[1],
		Document:  document,
		BirthDate: record[3],
		BetNumber: betNumber,
	}, nil
}
