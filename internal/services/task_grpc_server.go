package services

import (
	"context"
	"fmt"

	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/pb/taskpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TaskServer is the server that provides task services
type TaskServer struct {
	taskpb.UnimplementedTaskServiceServer // implement an interface with private method adding anonymous field for private method
	createTaskUseCase                     contracts.CreateTaskUseCase
	findTaskUseCase                       contracts.FindTaskUseCase
	listTaskUseCase                       contracts.ListTaskUseCase
	loginUseCase                          contracts.LoginUseCase
}

//NewTaskServer returns a new TaskServer
func NewTaskServer(createTaskUseCase contracts.CreateTaskUseCase,
	findTaskUseCase contracts.FindTaskUseCase,
	listTaskUseCase contracts.ListTaskUseCase,
	loginUseCase contracts.LoginUseCase) *TaskServer {

	return &TaskServer{
		createTaskUseCase: createTaskUseCase,
		findTaskUseCase:   findTaskUseCase,
		listTaskUseCase:   listTaskUseCase,
		loginUseCase:      loginUseCase,
	}
}

//CreateTask is a unary RPC to create a new task
func (t *TaskServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	task := req.Task
	fmt.Printf("receive a create-task with: %s", task.Summary)

	createTaskInput := &input.CreateTaskInput{
		Summary:     task.Summary,
		PerformedAt: task.PerformedAt.AsTime(),
		UserID:      task.UserId,
	}

	output, err := t.createTaskUseCase.Execute(ctx, createTaskInput)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if output == nil {
		return nil, status.Error(codes.Internal, "could not create task, try again")
	}

	res := &taskpb.CreateTaskResponse{
		Id: output.TaskID,
	}

	return res, nil
}

//Find a task by TaskID
func (t *TaskServer) FindTask(ctx context.Context, req *taskpb.FindTaskRequest) (*taskpb.FindTaskResponse, error) {
	fmt.Printf("receive a create-task request with: %s", req.Id)

	output, err := t.findTaskUseCase.Execute(ctx, req.UserId, req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error find task at database: %v", err)
	}

	if output == nil {
		return nil, status.Error(codes.Internal, "could not find task, try again")
	}

	res := &taskpb.FindTaskResponse{Task: &taskpb.Task{
		Id:          output.Task.ID,
		Summary:     output.Task.Summary,
		PerformedAt: timestamppb.New(output.Task.PerformedAt),
		UserId:      output.Task.UserID,
	}}

	return res, nil
}

//Perform List Tasks
func (t *TaskServer) ListTasks(ctx context.Context, req *taskpb.ListTaskRequest) (*taskpb.ListTaskResponse, error) {

	output, err := t.listTaskUseCase.Execute(ctx, req.UserId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error listing tasks at database: %v", err)
	}

	if output == nil {
		return nil, status.Error(codes.Internal, "could not list tasks, try again")
	}

	var taskspb []*taskpb.Task

	for _, taskEntity := range output.Tasks {
		t := &taskpb.Task{
			Id:          taskEntity.ID,
			UserId:      taskEntity.UserID,
			PerformedAt: timestamppb.New(taskEntity.PerformedAt),
			Summary:     taskEntity.Summary,
		}
		taskspb = append(taskspb, t)
	}

	res := &taskpb.ListTaskResponse{
		Tasks: taskspb,
	}

	return res, nil
}

func (t *TaskServer) Login(ctx context.Context, req *taskpb.LoginRequest) (*taskpb.LoginResponse, error) {

	output, err := t.loginUseCase.Execute(ctx, req.Email, req.Password)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error when try to login: %v", err)
	}

	if output == nil {
		return nil, status.Error(codes.Internal, "could not login, try again")
	}

	userID := ""

	if output.User != nil {
		userID = output.User.ID
	}

	res := &taskpb.LoginResponse{
		UserId: userID,
	}

	return res, nil
}
