package services_test

import (
	"testing"
)

func TestServerCreateTask(t *testing.T) {

	//init Inversion of control DI
	//dependencies := di.NewBuild()

	//when invalid userId should return an error

	//when database error should return an error

	//when create a test for manager user(role manager) trying to create a task(not allowed)

	//create task success

	/*
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
	*/
}
