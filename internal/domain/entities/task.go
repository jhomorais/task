package entities

import "time"

type Task struct {
	ID         string `gorm:"id"`
	Summary    string
	RealizedAt time.Time
	UserID     string `gorm:"index"`
	User       User
}
