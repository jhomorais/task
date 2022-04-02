package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerCreateTask(t *testing.T) {
	t.Parallel()

	//init dependencies, inversion of control DI
	dependencies := di.NewBuild()

	taskQueueMock := &mocks.RabbitMQ{}
	taskQueueMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	createTaskUseCase := usecases.NewCreateTaskUseCase(dependencies.Repositories.UserRepository,
		dependencies.Repositories.TaskRepository, taskQueueMock)

	userTec, err := dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "technician@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, userTec)

	userManager, err := dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "manager@gmail.com")
	assert.NoError(t, err)
	assert.NotNil(t, userManager)

	t.Run("when invalid userId should return an error", func(t *testing.T) {
		t.Parallel()

		taskNew := sample.NewTaskPb("invalid-user-id")

		req := &taskpb.CreateTaskRequest{
			Task: taskNew,
		}

		ctx := context.Background()
		server := services.NewTaskServer(createTaskUseCase,
			dependencies.Usecases.FindTaskUseCase,
			dependencies.Usecases.ListTaskUseCase)

		res, err := server.CreateTask(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.NotEqual(t, st.Code(), codes.OK)
	})

	t.Run("when database error", func(t *testing.T) {
		t.Parallel()

		errDatabase := errors.New("databases error")

		taskRepositoryErrorMock := &mocks.TaskRepository{}
		taskRepositoryErrorMock.On("CreateTask", mock.Anything, mock.Anything).Return(errDatabase)

		createTaskUseCase := usecases.NewCreateTaskUseCase(dependencies.Repositories.UserRepository,
			taskRepositoryErrorMock, dependencies.TaskQueue)

		taskNew := sample.NewTaskPb(userTec.ID)

		req := &taskpb.CreateTaskRequest{
			Task: taskNew,
		}

		ctx := context.Background()
		server := services.NewTaskServer(createTaskUseCase,
			dependencies.Usecases.FindTaskUseCase,
			dependencies.Usecases.ListTaskUseCase)

		res, err := server.CreateTask(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.NotEqual(t, st.Code(), codes.OK)
	})

	t.Run("when create a task with an user(role manager) should return an error", func(t *testing.T) {
		t.Parallel()

		taskNew := sample.NewTaskPb(userManager.ID)

		req := &taskpb.CreateTaskRequest{
			Task: taskNew,
		}

		ctx := context.Background()
		server := services.NewTaskServer(createTaskUseCase,
			dependencies.Usecases.FindTaskUseCase,
			dependencies.Usecases.ListTaskUseCase)

		res, err := server.CreateTask(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.NotEqual(t, st.Code(), codes.OK)
	})

	t.Run("when success create task", func(t *testing.T) {
		t.Parallel()

		taskNew := sample.NewTaskPb(userTec.ID)

		req := &taskpb.CreateTaskRequest{
			Task: taskNew,
		}

		ctx := context.Background()
		server := services.NewTaskServer(createTaskUseCase,
			dependencies.Usecases.FindTaskUseCase,
			dependencies.Usecases.ListTaskUseCase)

		res, err := server.CreateTask(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, res)

		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, st.Code(), codes.OK)
		assert.NotEmpty(t, res.Id)
	})

}
