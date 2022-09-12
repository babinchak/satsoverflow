package models

import (
	"time"
)

type User struct {
	ID       uint
	Username string `gorm:"primaryKey"`
	Email    string
	// Age          uint8
	Password      string
	Balance       uint
	TwitterHandle string
	// MemberNumber sql.NullString
	// ActivatedAt  sql.NullTime
	CreatedAt time.Time
	UpdatedAt time.Time
}
