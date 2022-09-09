package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           uint
	Username     string
	Email        string
	Age          uint8
	Password     string
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
