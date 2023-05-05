package model

import "time"

type Transaction struct {
	CitizenID string    `json:"citizen_id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Total     int       `json:"total"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTransaction(citizenID, from, to string, total int, timestamp time.Time) Transaction {
	return Transaction{
		CitizenID: citizenID,
		From:      from,
		To:        to,
		Total:     total,
		Timestamp: timestamp,
	}
}
