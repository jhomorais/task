package repositories

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, entity *entities.User) error
	FindById(ctx context.Context, id string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, enitity *entities.User) error
}
