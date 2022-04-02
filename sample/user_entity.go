package sample

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fgmaia/task/internal/domain/entities"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func NewUserEntity() *entities.User {

	randomRole := RandomStringFromSet(entities.ROLE_TECHNICIAN.String(), entities.ROLE_MANAGER.String())
	return NewUserEntityRole(entities.Role(randomRole))
}

func NewUserEntityRole(role entities.Role) *entities.User {

	user := &entities.User{
		ID:       RandomID(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 8),
		Role:     role,
	}

	return user
}
