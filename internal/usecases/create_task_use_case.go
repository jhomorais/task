package usecases

import (
	"context"
	"fmt"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/infra/queue"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/internal/usecases/ports/output"
	"github.com/fgmaia/task/internal/usecases/validator"
	"github.com/google/uuid"
)

type createTaskUseCase struct {
	userRepository repositories.UserRepository
	taskRepository repositories.TaskRepository
	taskQueue      queue.RabbitMQ
}

func NewCreateTaskUseCase(userRepository repositories.UserRepository,
	taskRepository repositories.TaskRepository,
	taskQueue queue.RabbitMQ) contracts.CreateTaskUseCase {

	return &createTaskUseCase{
		userRepository: userRepository,
		taskRepository: taskRepository,
		taskQueue:      taskQueue,
	}
}

func (c *createTaskUseCase) Execute(ctx context.Context, createTask *input.CreateTaskInput) (*output.CreateTaskOutput, error) {

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate a new task ID: %v", err)
	}

	if err := validator.ValidateUUId(createTask.UserID, true, "userId"); err != nil {
		return nil, err
	}

	userEntity, err := c.userRepository.FindById(ctx, createTask.UserID)
	if err != nil {
		return nil, fmt.Errorf("erro when try to find user: %v", err)
	}

	if userEntity == nil || userEntity.ID == "" {
		return nil, fmt.Errorf("userID not found")
	}

	if userEntity.Role != entities.ROLE_TECHNICIAN {
		return nil, fmt.Errorf("invalid user role only technicians can create tasks : %s", userEntity.ID)
	}

	taskEntity := &entities.Task{
		ID:          id.String(),
		UserID:      createTask.UserID,
		Summary:     createTask.Summary,
		PerformedAt: createTask.PerformedAt,
	}

	err = c.taskRepository.CreateTask(ctx, taskEntity)

	if err != nil {
		return nil, fmt.Errorf("cannot save task at database: %v", err)
	}

	msg := fmt.Sprintf("User: %s performed taskId: %s summary: %s", userEntity.Email, taskEntity.ID, taskEntity.Summary)
	err = c.taskQueue.Publish(ctx, []byte(msg))
	if err != nil {
		fmt.Println(err) //only log error
	}

	createTaskOutput := &output.CreateTaskOutput{
		TaskID: taskEntity.ID,
	}

	return createTaskOutput, nil
}
