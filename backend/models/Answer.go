package models

import "time"

type Answer struct {
	ID         uint
	Body       string
	Bounty     uint
	QuestionID uint
	AnswerName *string
	Answer     User `gorm:"foreignKey:AnswerName"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
