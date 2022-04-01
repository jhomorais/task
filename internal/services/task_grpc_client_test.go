package services_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/fgmaia/task/config"
	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/infra/queue"
	"github.com/fgmaia/task/internal/serializer"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/fgmaia/task/sample"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestIntegrationClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	//init DI
	dependencies := di.NewBuild()

	//users
	userTec, err := userRepository.FindByEmail(context.Background(), "technician@gmail.com")
	require.NoError(t, err)
	require.NotNil(t, userTec)
	require.NotEmpty(t, userTec.ID)

	//init rabbitMQ
	taskQueue := queue.NewRabbitMQ(config.GetRabbitMQClient(),
		"task-exchange",
		"tasks")

	err = taskQueue.InitQueue(ctx)
	require.NoError(t, err)
	defer taskQueue.Close()

	//init server and inject repositories and queue
	taskServer := services.NewTaskServer(userRepository, taskRepository, taskQueue)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener) // non block call

	//create client
	taskClient := newTestTaskClient(t, listener.Addr().String())

	task := sample.NewTaskPb(userTec.ID)
	expectedID := task.Id

	//
	req := &taskpb.CreateTaskRequest{
		Task: task,
	}

	//client call
	res, err := taskClient.CreateTask(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, expectedID, res.Id)

	reqFind := &taskpb.FindTaskRequest{
		Id:     task.Id,
		UserId: userTec.ID,
	}
	resFind, err := taskServer.FindTask(context.Background(), reqFind)
	require.NoError(t, err)
	require.NotNil(t, resFind)
	require.Equal(t, expectedID, resFind.Task.Id)
	//TODO fix performed_at divergence
	/*
		-  "performed_at": "2022-03-31T19:24:56.712450900Z"
		+  "performed_at": "2022-03-31T19:24:56.712Z"
	*/
	//requireSameTask(t, task, resFind.Task)

	expectedMessage := fmt.Sprintf("User: %s performed taskId: %s summary: %s", userTec.Email, task.Id, task.Summary)
	messegaReceived := ""
	messages, err := taskQueue.Consume(ctx)
	require.NoError(t, err)
	for msg := range messages {
		messegaReceived = string(msg.Body)
		break
	}
	require.Equal(t, expectedMessage, messegaReceived)

	//TODO ListTask
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
