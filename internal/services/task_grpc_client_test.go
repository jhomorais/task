package services_test

import (
	"context"
	"net"
	"testing"

	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/serializer"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestGrpcClient(t *testing.T) {
	t.Parallel()

	//init
	dependencies := di.NewBuild()

	taskQueueMock := &mocks.RabbitMQ{}
	taskQueueMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	createTaskUseCase := usecases.NewCreateTaskUseCase(dependencies.Repositories.UserRepository,
		dependencies.Repositories.TaskRepository, taskQueueMock)

	userTec, err := dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "technician@gmail.com")
	require.NoError(t, err)
	require.NotNil(t, userTec)
	require.NotEmpty(t, userTec.ID)

	userManager, err := dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "manager@gmail.com")
	require.NoError(t, err)
	require.NotNil(t, userManager)
	require.NotEmpty(t, userManager.ID)

	//init server and inject repositories and queue
	taskServer := services.NewTaskServer(createTaskUseCase,
		dependencies.Usecases.FindTaskUseCase,
		dependencies.Usecases.ListTaskUseCase)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non block call

	taskClient := newTestTaskClient(t, listener.Addr().String())

	t.Run("when try to create a task with invalid user", func(t *testing.T) {
		t.Parallel()

		task := sample.NewTaskPb("invalid-user-id")

		req := &taskpb.CreateTaskRequest{
			Task: task,
		}

		res, err := taskClient.CreateTask(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("when try to create a task with valid user with an invalid role", func(t *testing.T) {
		t.Parallel()

		task := sample.NewTaskPb(userManager.ID)

		req := &taskpb.CreateTaskRequest{
			Task: task,
		}

		res, err := taskClient.CreateTask(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("when success create and find a task", func(t *testing.T) {
		t.Parallel()

		task := sample.NewTaskPb(userTec.ID)

		req := &taskpb.CreateTaskRequest{
			Task: task,
		}

		res, err := taskClient.CreateTask(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res.Id)

		//
		reqFindTask := &taskpb.FindTaskRequest{
			Id:     res.Id,
			UserId: task.UserId,
		}

		resFindTask, err := taskClient.FindTask(context.Background(), reqFindTask)
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, resFindTask.Task.Id)
		require.Equal(t, resFindTask.Task.Id, reqFindTask.Id)
	})

	t.Run("when list tasks with invalid user role", func(t *testing.T) {
		t.Parallel()

		req := &taskpb.ListTaskRequest{
			UserId: userTec.ID,
		}

		res, err := taskClient.ListTasks(context.Background(), req)
		require.Error(t, err)
		require.Nil(t, res)
	})

	t.Run("when success list tasks", func(t *testing.T) {
		t.Parallel()

		req := &taskpb.ListTaskRequest{
			UserId: userManager.ID,
		}

		res, err := taskClient.ListTasks(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, res)
	})

}

func newTestTaskClient(t *testing.T, serverAddress string) taskpb.TaskServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return taskpb.NewTaskServiceClient(conn)
}

func requireSameTask(t *testing.T, task1 *taskpb.Task, task2 *taskpb.Task) {
	json1, err := serializer.ProtobufToJSON(task1)
	require.NoError(t, err)

	json2, err := serializer.ProtobufToJSON(task2)
	require.NoError(t, err)

	require.Equal(t, json1, json2)
}
