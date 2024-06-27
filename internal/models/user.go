package models

import (
	"github.com/gocql/gocql"
	"time"
)

type User struct {
	ID        gocql.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}
