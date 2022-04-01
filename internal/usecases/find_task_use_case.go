package usecases

import (
	"context"
	"fmt"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/input"
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

func (c *findTaskUseCase) Execute(ctx context.Context, findTask *input.FindTaskInput) (*output.FindTaskOutput, error) {

	if err := validator.ValidateUUId(findTask.TaskID, true, "taskId"); err != nil {
		return nil, err
	}

	if err := validator.ValidateUUId(findTask.UserID, true, "userId"); err != nil {
		return nil, err
	}

	userEntity, err := c.userRepository.FindById(ctx, findTask.UserID)
	if err != nil {
		return nil, fmt.Errorf("erro when try to find user: %v", err)
	}

	if userEntity == nil || userEntity.ID == "" {
		return nil, fmt.Errorf("userID not found")
	}

	taskEntity, err := c.taskRepository.FindTask(ctx, findTask.TaskID)

	if userEntity.Role == entities.ROLE_TECHNICIAN {
		if userEntity.ID != taskEntity.User.ID {
			return nil, fmt.Errorf("error this task do not belong to the user")
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
