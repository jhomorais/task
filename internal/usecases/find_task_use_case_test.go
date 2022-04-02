package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFindTask(t *testing.T) {
	t.Parallel()

	userTec1 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
	task1 := sample.NewTaskEntityWithUser(*userTec1)

	taskRepositoryMock := &mocks.TaskRepository{}
	taskRepositoryMock.On("FindTask", mock.Anything, task1.ID).Return(task1, nil)

	t.Run("when taskId is invalid should return an error", func(t *testing.T) {
		t.Parallel()

		taskNoId := ""
		taskInvalidId := "invalid-uuid"

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(&userTec1, nil)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userTec1.ID, taskNoId)
		require.Error(t, err)
		require.Nil(t, output)

		output, err = findTaskUseCase.Execute(context.Background(), userTec1.ID, taskInvalidId)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when userId is invalid should return an error", func(t *testing.T) {
		t.Parallel()

		userNoId := ""
		userInvalidId := "invalid-uuid"

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(&userTec1, nil)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userNoId, task1.ID)
		require.Error(t, err)
		require.Nil(t, output)

		output, err = findTaskUseCase.Execute(context.Background(), userInvalidId, task1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when not found user should return an error", func(t *testing.T) {
		t.Parallel()

		errUserNotFound := errors.New("erro when try to find user")

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(nil, errUserNotFound)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userTec1.ID, task1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when finding task error", func(t *testing.T) {
		t.Parallel()

		errTaskNotFound := errors.New("erro when try to find task")

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(userTec1, nil)

		taskRepositoryMock := &mocks.TaskRepository{}
		taskRepositoryMock.On("FindTask", mock.Anything, task1.ID).Return(nil, errTaskNotFound)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userTec1.ID, task1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when user role is equals to TECHNICIAN and user is not the owner of it", func(t *testing.T) {
		t.Parallel()

		userTec2 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec2.ID).Return(userTec2, nil)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userTec2.ID, task1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when successfully FindTask", func(t *testing.T) {
		t.Parallel()

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(userTec1, nil)

		findTaskUseCase := usecases.NewFindTaskUseCase(userRepositoryMock, taskRepositoryMock)

		output, err := findTaskUseCase.Execute(context.Background(), userTec1.ID, task1.ID)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

}
