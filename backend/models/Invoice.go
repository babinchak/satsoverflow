package models

import "time"

// Invoice is for users adding funds to their account
type Invoice struct {
	ID        uint
	Hash      string
	Username  *string
	User      User `gorm:"foreignKey:Username"`
	Paid      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
