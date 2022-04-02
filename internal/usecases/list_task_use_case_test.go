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

func TestListTask(t *testing.T) {
	t.Parallel()

	userTec1 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
	userTec2 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
	userTec3 := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)

	var listTask = []*entities.Task{
		sample.NewTaskEntityWithUser(*userTec1),
		sample.NewTaskEntityWithUser(*userTec2),
		sample.NewTaskEntityWithUser(*userTec1),
		sample.NewTaskEntityWithUser(*userTec2),
		sample.NewTaskEntityWithUser(*userTec3),
	}

	taskRepositoryMock := &mocks.TaskRepository{}
	taskRepositoryMock.On("ListTask", mock.Anything).Return(listTask, nil)

	t.Run("when userId is invalid should return an error", func(t *testing.T) {
		t.Parallel()

		userNoId := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
		userNoId.ID = ""

		userInvalidId := sample.NewUserEntityRole(entities.ROLE_TECHNICIAN)
		userInvalidId.ID = "invalid-uuid"

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", mock.Anything, userTec1.ID).Return(&userTec1, nil)

		listTaskUseCase := usecases.NewListTaskUseCase(userRepositoryMock,
			taskRepositoryMock)

		output, err := listTaskUseCase.Execute(context.Background(), userNoId.ID)
		require.Error(t, err)
		require.Nil(t, output)

		output, err = listTaskUseCase.Execute(context.Background(), userInvalidId.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when not found user should return an error", func(t *testing.T) {
		t.Parallel()

		errUserNotFound := errors.New("erro when try to find user")

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", context.Background(), userTec1.ID).Return(nil, errUserNotFound)

		listTaskUseCase := usecases.NewListTaskUseCase(userRepositoryMock,
			taskRepositoryMock)

		output, err := listTaskUseCase.Execute(context.Background(), userTec1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when user role is differente of MANAGER should return an error", func(t *testing.T) {
		t.Parallel()

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", context.Background(), userTec1.ID).Return(userTec1, nil)

		listTaskUseCase := usecases.NewListTaskUseCase(userRepositoryMock,
			taskRepositoryMock)

		output, err := listTaskUseCase.Execute(context.Background(), userTec1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when ListTask return an error", func(t *testing.T) {
		t.Parallel()

		errDatabase := errors.New("databases error")
		userManager := sample.NewUserEntityRole(entities.ROLE_MANAGER)

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", context.Background(), userTec1.ID).Return(userManager, nil)

		taskRepositoryErrorMock := &mocks.TaskRepository{}
		taskRepositoryErrorMock.On("ListTask", mock.Anything).Return(nil, errDatabase)

		listTaskUseCase := usecases.NewListTaskUseCase(userRepositoryMock,
			taskRepositoryErrorMock)

		output, err := listTaskUseCase.Execute(context.Background(), userTec1.ID)
		require.Error(t, err)
		require.Nil(t, output)
	})

	t.Run("when successfully ListTask", func(t *testing.T) {
		t.Parallel()

		userManager := sample.NewUserEntityRole(entities.ROLE_MANAGER)

		userRepositoryMock := &mocks.UserRepository{}
		userRepositoryMock.On("FindById", context.Background(), userTec1.ID).Return(userManager, nil)

		listTaskUseCase := usecases.NewListTaskUseCase(userRepositoryMock,
			taskRepositoryMock)

		output, err := listTaskUseCase.Execute(context.Background(), userTec1.ID)
		require.NoError(t, err)
		require.NotNil(t, output)
	})

}
