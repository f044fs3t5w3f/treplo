package models

import "time"

type User struct {
	ID        int64
	FirstName string
	LastName  string
	CreatedAt time.Time
}
