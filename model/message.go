package model

import "time"

type Message struct {
	CitizenID string    `json:"citizen_id"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewMessage(citizenID, message string, timestamp time.Time) Message {
	return Message{
		CitizenID: citizenID,
		Message:   message,
		Timestamp: timestamp,
	}
}
