package services_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/serializer"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/mocks"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var (
	dependencies  *di.DenpencyBuild
	taskQueueMock *mocks.RabbitMQ
	userTec       *entities.User
	userManager   *entities.User
	taskServer    *services.TaskServer
	listenerAddr  string
)

func initDependencies(t *testing.T) {
	if dependencies != nil {
		return
	}
	var err error

	dependencies = di.NewBuild()

	imageStore := services.NewDiskImageStore(".")

	taskQueueMock = &mocks.RabbitMQ{}
	taskQueueMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	userTec, err = dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "technician@gmail.com")
	require.NoError(t, err)
	require.NotNil(t, userTec)
	require.NotEmpty(t, userTec.ID)

	userManager, err = dependencies.Repositories.UserRepository.FindByEmail(context.Background(), "manager@gmail.com")
	require.NoError(t, err)
	require.NotNil(t, userManager)
	require.NotEmpty(t, userManager.ID)

	createTaskUseCase := usecases.NewCreateTaskUseCase(dependencies.Repositories.UserRepository,
		dependencies.Repositories.TaskRepository, taskQueueMock)

	taskServer = services.NewTaskServer(createTaskUseCase,
		dependencies.Usecases.FindTaskUseCase,
		dependencies.Usecases.ListTaskUseCase,
		dependencies.Usecases.LoginUseCase,
		imageStore)
}

func TestClientCreateTask(t *testing.T) {
	t.Parallel()

	initDependencies(t)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	listenerAddr = listener.Addr().String()
	go grpcServer.Serve(listener) // non block call

	taskClient := newTestTaskClient(t, listenerAddr)

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
}

func TestClientListTasks(t *testing.T) {
	t.Parallel()

	initDependencies(t)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	listenerAddr = listener.Addr().String()
	go grpcServer.Serve(listener) // non block call

	taskClient := newTestTaskClient(t, listenerAddr)

	t.Run("when list tasks with invalid user role", func(t *testing.T) {
		t.Parallel()

		req := &taskpb.ListTaskRequest{
			UserId: userTec.ID,
		}

		stream, err := taskClient.ListTasks(context.Background(), req)
		require.NoError(t, err)

		_, err = stream.Recv()
		require.Error(t, err)
	})

	t.Run("when success list tasks", func(t *testing.T) {
		t.Parallel()

		req := &taskpb.ListTaskRequest{
			UserId: userManager.ID,
		}

		stream, err := taskClient.ListTasks(context.Background(), req)
		require.NoError(t, err)

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			require.NoError(t, err)
			require.NotEmpty(t, res.Task.Id)
		}
	})

}

func TestClienteUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../.."

	initDependencies(t)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	listenerAddr = listener.Addr().String()
	go grpcServer.Serve(listener) // non block call

	taskClient := newTestTaskClient(t, listenerAddr)

	taskInput := &input.CreateTaskInput{
		Summary:     "teste uploda image client",
		PerformedAt: time.Now(),
		UserID:      userTec.ID,
	}

	taskOutput, err := dependencies.Usecases.CreateTaskUseCase.Execute(context.Background(), taskInput)
	require.NoError(t, err)
	require.NotNil(t, taskOutput)
	require.NotEmpty(t, taskOutput.TaskID)

	imagePath := fmt.Sprintf("%s/how_to_test_console.png", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := taskClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &taskpb.UploadImageRequest{
		UserId: userTec.ID,
		Data: &taskpb.UploadImageRequest_Info{
			Info: &taskpb.ImageInfo{
				TaskId:    taskOutput.TaskID,
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n

		req := &taskpb.UploadImageRequest{
			Data: &taskpb.UploadImageRequest_ChunckData{
				ChunckData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s%s", res.GetId(), imageType)
	require.FileExists(t, savedImagePath)
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
