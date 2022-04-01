package input

import "time"

type CreateTaskInput struct {
	Summary     string
	PerformedAt time.Time
	UserID      string
}
