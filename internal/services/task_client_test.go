package services_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/fgmaia/task/internal/infra/config"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/serializer"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestClientCreateTask(t *testing.T) {
	t.Parallel()

	taskServer, serverAddress := startTestTaskServer(t)
	taskClient := newTestTaskClient(t, serverAddress)

	task := sample.NewTaskPb()
	expectedID := task.Id

	req := &taskpb.CreateTaskRequest{
		Task: task,
	}

	res, err := taskClient.CreateTask(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	reqFind := &taskpb.FindTaskRequest{
		Id: task.Id,
	}
	resFind, err := taskServer.FindTask(context.Background(), reqFind)
	require.NoError(t, err)
	require.NotNil(t, resFind)
	require.Equal(t, expectedID, resFind.Task.Id)
	requireSameTask(t, task, resFind.Task)
}

func startTestTaskServer(t *testing.T) (*services.TaskServer, string) {
	//could mock db, but will be hard to test FindTask
	///taskRepository := &mocks.TaskRepository{}
	//taskRepository.On("CreateTask", mock.Anything, mock.Anything).Return(nil)
	db := initGormMysqlDB(t)
	taskRepository := repositories.NewTaskRepository(db)

	taskServer := services.NewTaskServer(taskRepository)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non block call

	return taskServer, listener.Addr().String()
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

func initGormMysqlDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@%s", config.GetMysqlUser(), config.GetMysqlPassword(), config.GetMysqlConnectionString())
	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	mysqlDb.AutoMigrate()
	return mysqlDb
}
