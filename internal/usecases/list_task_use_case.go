package usecases

import (
	"context"
	"fmt"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/validator"
)

type listTaskUseCase struct {
	userRepository repositories.UserRepository
	taskRepository repositories.TaskRepository
}

func NewListTaskUseCase(userRepository repositories.UserRepository,
	taskRepository repositories.TaskRepository) contracts.ListTaskUseCase {

	return &listTaskUseCase{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (l *listTaskUseCase) Execute(ctx context.Context, userID string, found func(task *entities.Task) error) error {

	if err := validator.ValidateUUId(userID, true, "userId"); err != nil {
		return err
	}

	userEntity, err := l.userRepository.FindById(ctx, userID)
	if err != nil {
		return fmt.Errorf("erro when try to find user: %v", err)
	}

	if userEntity == nil || userEntity.ID == "" {
		return fmt.Errorf("userId is not valid: %s", userEntity.ID)
	}

	if userEntity.Role != entities.ROLE_MANAGER {
		return fmt.Errorf("invalid user role, only managers can list tasks")
	}

	err = l.taskRepository.ListTask(ctx, found)

	if err != nil {
		return fmt.Errorf("error when list tasks on database: %v", err)
	}

	return nil
}
