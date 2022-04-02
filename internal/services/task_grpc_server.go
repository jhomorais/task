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
}

//NewTaskServer returns a new TaskServer
func NewTaskServer(createTaskUseCase contracts.CreateTaskUseCase,
	findTaskUseCase contracts.FindTaskUseCase,
	listTaskUseCase contracts.ListTaskUseCase) *TaskServer {

	return &TaskServer{
		createTaskUseCase: createTaskUseCase,
		findTaskUseCase:   findTaskUseCase,
		listTaskUseCase:   listTaskUseCase,
	}
}

//CreateTask is a unary RPC to create a new task
func (t *TaskServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	task := req.Task
	fmt.Printf("receive a create-task request with: %s", task.Id)

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
