package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

//CreateTask test
func TestCreateTask(t *testing.T) {
	t.Parallel()

	taskQueueMock := &mocks.RabbitMQ{}
	taskQueueMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	userTec1 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
	task1 := sample.NewTaskEntityWithUser(*userTec1)

	taskRepositoryMock := &mocks.TaskRepository{}
	taskRepositoryMock.On("CreateTask", mock.Anything, mock.Anything).Return(nil)

	createTaskInput := &input.CreateTaskInput{
		Summary:     task1.Summary,
		PerformedAt: task1.PerformedAt,
		UserID:      userTec1.ID,
	}

	t.Run("when userId is empty should return an error", func(t *testing.T) {
		t.Parallel()

		userNoId := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
		userNoId.ID = ""

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userNoId.ID).Return(userNoId, nil)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryMock, taskQueueMock)

		createTaskInput := &input.CreateTaskInput{
			Summary:     task1.Summary,
			PerformedAt: task1.PerformedAt,
			UserID:      userNoId.ID,
		}

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when userId is invalid should return an error", func(t *testing.T) {
		t.Parallel()

		userInvalidId := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
		userInvalidId.ID = "invalid-uuid"

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userInvalidId.ID).Return(userInvalidId, nil)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryMock, taskQueueMock)

		createTaskInput := &input.CreateTaskInput{
			Summary:     task1.Summary,
			PerformedAt: task1.PerformedAt,
			UserID:      userInvalidId.ID,
		}

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when not found user should return an error", func(t *testing.T) {
		t.Parallel()

		errUserNotFound := errors.New("erro when try to find user")

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(nil, errUserNotFound)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryMock, taskQueueMock)

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when user role is different of TECHNICIAN should return an error", func(t *testing.T) {
		t.Parallel()

		userManager := sample.NewUserEntityRole(entities.ROLE_MANAGER)

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userManager.ID).Return(userManager, nil)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryMock, taskQueueMock)

		createTaskInput := &input.CreateTaskInput{
			Summary:     task1.Summary,
			PerformedAt: task1.PerformedAt,
			UserID:      userManager.ID,
		}

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when CreateTask db store returns an error", func(t *testing.T) {
		t.Parallel()

		errDatabase := errors.New("databases error")

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(userTec1, nil)

		taskRepositoryErrorMock := &mocks.TaskRepository{}
		taskRepositoryErrorMock.On("CreateTask", mock.Anything, mock.Anything).Return(errDatabase)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryErrorMock, taskQueueMock)

		createTaskInput := &input.CreateTaskInput{
			Summary:     task1.Summary,
			PerformedAt: task1.PerformedAt,
			UserID:      userTec1.ID,
		}

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when successfully CreateTask", func(t *testing.T) {
		t.Parallel()

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(userTec1, nil)

		createTaskUseCase := usecases.NewCreateTaskUseCase(userRepositoryMock,
			taskRepositoryMock, taskQueueMock)

		createTaskInput := &input.CreateTaskInput{
			Summary:     task1.Summary,
			PerformedAt: task1.PerformedAt,
			UserID:      userTec1.ID,
		}

		output, err := createTaskUseCase.Execute(context.Background(), createTaskInput)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

}
