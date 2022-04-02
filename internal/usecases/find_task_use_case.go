package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/output"
	"github.com/fgmaia/task/internal/usecases/validator"
)

type findTaskUseCase struct {
	userRepository repositories.UserRepository
	taskRepository repositories.TaskRepository
}

func NewFindTaskUseCase(userRepository repositories.UserRepository,
	taskRepository repositories.TaskRepository) contracts.FindTaskUseCase {

	return &findTaskUseCase{
		userRepository: userRepository,
		taskRepository: taskRepository,
	}
}

func (c *findTaskUseCase) Execute(ctx context.Context, userID string, taskID string) (*output.FindTaskOutput, error) {

	if err := validator.ValidateUUId(taskID, true, "taskId"); err != nil {
		return nil, err
	}

	if err := validator.ValidateUUId(userID, true, "userId"); err != nil {
		return nil, err
	}

	userEntity, err := c.userRepository.FindById(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("erro when try to find user: %v", err)
	}

	if userEntity == nil || userEntity.ID == "" {
		return nil, fmt.Errorf("userID not found")
	}

	taskEntity, err := c.taskRepository.FindTask(ctx, taskID)

	if err != nil {
		return nil, fmt.Errorf("error when try to find a task: %w", err)
	}

	if userEntity.Role == entities.ROLE_TECHNICIAN {
		if userEntity.ID != taskEntity.User.ID {
			return nil, errors.New("error this task do not belong to the user")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error find task at database: %v", err)
	}

	output := &output.FindTaskOutput{
		Task: taskEntity,
	}

	return output, nil
}
