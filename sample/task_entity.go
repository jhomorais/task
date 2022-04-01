package sample

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fgmaia/task/internal/domain/entities"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func NewTaskEntity() *entities.Task {

	task := &entities.Task{
		ID:          randomID(),
		Summary:     randomSummary(),
		PerformedAt: time.Now(),
		UserID:      randomID(),
	}

	return task
}
