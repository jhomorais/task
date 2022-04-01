package entities

import "time"

type Task struct {
	ID          string `gorm:"id"`
	Summary     string `gorm:"size:2500"`
	PerformedAt time.Time
	UserID      string `gorm:"index"`
	User        User
}
