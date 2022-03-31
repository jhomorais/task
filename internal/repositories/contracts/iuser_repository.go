package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, entity *entities.User) error
}
