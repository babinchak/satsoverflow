package models

import "time"

type Answer struct {
	ID         uint
	Body       string
	Bounty     uint
	QuestionID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
