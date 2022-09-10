package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID           uint
	Username     string `gorm:"primaryKey"`
	Email        string
	Age          uint8
	Password     string
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
