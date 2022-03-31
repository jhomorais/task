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

	randomRole := randomStringFromSet(entities.ROLE_TECHNICIAN.String(), entities.ROLE_MANAGER.String())

	user := &entities.User{
		ID:       randomID(),
		Email:    gofakeit.Email(),
		Password: gofakeit.Password(true, true, true, true, true, 8),
		Role:     entities.Role(randomRole),
	}

	return user
}
