package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, entity *entities.Task) error
	FindTask(ctx context.Context, id string) (*entities.Task, error)
}
