package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/infra/queue"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//********* If prefer can mocks access to repository *************

//********* If prefer can mocks access to repository *************

func TestServerCreateTask(t *testing.T) {
	t.Parallel()

	//init DI
	dependencies := di.NewBuild()

	//mocks
	taskRepositoryErrorMock := &mocks.TaskRepository{}
	taskRepositoryErrorMock.On("CreateTask", mock.Anything, mock.Anything).Return(errors.New("databases error"))

	userRepositoryErrorMock := &mocks.UserRepository{}
	userRepositoryErrorMock.On("FindById", mock.Anything, mock.Anything).Return(userTech, nil)

	taskNoId := sample.NewTaskPb(user.ID)
	taskNoId.Id = ""

	taskInvalidId := sample.NewTaskPb(user.ID)
	taskInvalidId.Id = "invalid-uuid"

	//TODO create a test for manager user(role manager) trying to create a task(not allowed)

	testCases := []struct {
		name           string
		task           *taskpb.Task
		taskRepository contracts.TaskRepository
		userRepository contracts.UserRepository
		taskQueue      queue.RabbitMQ
		code           codes.Code
	}{
		{
			name:           "success_with_id",
			task:           sample.NewTaskPb(user.ID),
			taskRepository: taskRepositorySuccess,
			userRepository: userRepository,
			taskQueue:      taskQueue,
			code:           codes.OK,
		},
		{
			name:           "success_no_id",
			task:           taskNoId,
			taskRepository: taskRepositorySuccess,
			userRepository: userRepository,
			taskQueue:      taskQueue,
			code:           codes.OK,
		},
		{
			name:           "failure_invalid_id",
			task:           taskInvalidId,
			taskRepository: taskRepositorySuccess,
			userRepository: userRepository,
			taskQueue:      taskQueue,
			code:           codes.InvalidArgument,
		},
		{
			name:           "failure_on_database",
			task:           sample.NewTaskPb(user.ID),
			taskRepository: taskRepositoryError,
			userRepository: userRepository,
			taskQueue:      taskQueue,
			code:           codes.Internal,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &taskpb.CreateTaskRequest{
				Task: tc.task,
			}
			ctx := context.Background()
			server := services.NewTaskServer(tc.userRepository, tc.taskRepository, tc.taskQueue)
			res, err := server.CreateTask(ctx, req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.task.Id) > 0 {
					require.Equal(t, tc.task.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				_, ok := status.FromError(err)
				require.True(t, ok)
			}
		})

	}
}
