package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fgmaia/task/internal/repositories/contracts"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServerCreateTask(t *testing.T) {
	t.Parallel()

	taskRepositorySuccess := &mocks.TaskRepository{}
	taskRepositorySuccess.On("CreateTask", mock.Anything, mock.Anything).Return(nil)

	taskRepositoryError := &mocks.TaskRepository{}
	taskRepositoryError.On("CreateTask", mock.Anything, mock.Anything).Return(errors.New("databases error"))

	taskNoId := sample.NewTaskPb()
	taskNoId.Id = ""

	taskInvalidId := sample.NewTaskPb()
	taskInvalidId.Id = "invalid-uuid"

	testCases := []struct {
		name  string
		task  *taskpb.Task
		store contracts.TaskRepository
		code  codes.Code
	}{
		{
			name:  "success_with_id",
			task:  sample.NewTaskPb(),
			store: taskRepositorySuccess,
			code:  codes.OK,
		},
		{
			name:  "success_no_id",
			task:  taskNoId,
			store: taskRepositorySuccess,
			code:  codes.OK,
		},
		{
			name:  "failure_invalid_id",
			task:  taskInvalidId,
			store: taskRepositorySuccess,
			code:  codes.InvalidArgument,
		},
		{
			name:  "failure_on_database",
			task:  sample.NewTaskPb(),
			store: taskRepositoryError,
			code:  codes.Internal,
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
			server := services.NewTaskServer(tc.store)
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
