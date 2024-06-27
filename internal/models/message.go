package models

import (
	"github.com/gocql/gocql"
	"time"
)

type Message struct {
	ID        gocql.UUID `json:"id"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	Timestamp time.Time `json:"timestamp"`
	Content   string    `json:"content"`
}
