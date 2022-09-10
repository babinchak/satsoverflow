package models

import "time"

type Question struct {
	ID uint
	// Hash      string
	Title  string
	Body   string
	Bounty uint
	Paid   bool
	Hash   string
	// Asker     User
	CreatedAt time.Time
	UpdatedAt time.Time
}
