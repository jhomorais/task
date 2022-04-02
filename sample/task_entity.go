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
	user := NewUserEntityRole(entities.ROLE_TECHNICIAN)
	return NewTaskEntityWithUser(*user)
}

func NewTaskEntityWithUser(user entities.User) *entities.Task {
	task := &entities.Task{
		ID:          RandomID(),
		Summary:     RandomSummary(),
		PerformedAt: time.Now(),
		User:        user,
		UserID:      user.ID,
	}
	return task
}
