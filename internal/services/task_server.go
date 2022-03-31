package services

import (
	"context"
	"fmt"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories/contracts"
	"github.com/fgmaia/task/pb/taskpb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//TaskServer is the server that provides task services
type TaskServer struct {
	taskpb.UnimplementedTaskServiceServer // implement an interface with private method adding anonymous field for private method
	taskRepository                        contracts.TaskRepository
}

//NewTaskServer returns a new TaskServer
func NewTaskServer(taskRepository contracts.TaskRepository) *TaskServer {
	return &TaskServer{taskRepository: taskRepository}
}

//CreateTask is a unary RPC to create a new task
func (t *TaskServer) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	task := req.Task
	fmt.Printf("receive a create-task request with: %s", task.Id)

	if len(task.Id) > 0 {
		_, err := uuid.Parse(task.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "task ID is not valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new task ID: %v", err)
		}
		task.Id = id.String()
	}

	taskEntity := &entities.Task{ID: task.Id,
		Summary:    task.Summary,
		RealizedAt: task.PerformedAt.AsTime(),
		UserID:     task.UserId,
	}
	err := t.taskRepository.CreateTask(ctx, taskEntity)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save task at database: %v", err)
	}

	res := &taskpb.CreateTaskResponse{
		Id: taskEntity.ID,
	}

	return res, nil
}

func (t *TaskServer) FindTask(ctx context.Context, req *taskpb.FindTaskRequest) (*taskpb.FindTaskResponse, error) {
	fmt.Printf("receive a create-task request with: %s", req.Id)

	if len(req.Id) > 0 {
		_, err := uuid.Parse(req.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "task ID is not valid UUID: %v", err)
		}
	} else {
		return nil, status.Errorf(codes.Internal, "Id cannot be empty")
	}

	taskEntity, err := t.taskRepository.FindTask(ctx, req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot save task at database: %v", err)
	}

	res := &taskpb.FindTaskResponse{Task: &taskpb.Task{
		Id:          taskEntity.ID,
		Summary:     taskEntity.Summary,
		PerformedAt: timestamppb.New(taskEntity.RealizedAt),
		UserId:      taskEntity.UserID,
	}}

	return res, nil
}
