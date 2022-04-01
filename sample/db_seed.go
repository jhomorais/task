package sample

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DBSeed(db *gorm.DB) error {
	userRepository := repositories.NewUserRepository(db)

	err := createUser(userRepository, "technician@gmail.com", "tech", entities.ROLE_TECHNICIAN)

	if err != nil {
		return err
	}

	err = createUser(userRepository, "manager@gmail.com", "manager", entities.ROLE_MANAGER)

	if err != nil {
		return err
	}

	return nil
}

func createUser(userRepository repositories.UserRepository, email string, password string, role entities.Role) error {
	//
	ctx := context.Background()
	user, err := userRepository.FindByEmail(ctx, email)

	if err != nil {
		return err
	}

	id, err := uuid.NewUUID()

	if err != nil {
		return err
	}

	if user.ID == "" {
		user = &entities.User{
			ID:       id.String(),
			Email:    email,
			Password: password,
			Role:     role,
		}
		err := userRepository.CreateUser(context.Background(), user)
		if err != nil {
			return err
		}
	}

	return nil
}
