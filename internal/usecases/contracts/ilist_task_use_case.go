package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
)

type ListTaskUseCase interface {
	Execute(ctx context.Context, userID string, found func(task *entities.Task) error) error
}
