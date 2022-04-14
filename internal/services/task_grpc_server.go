package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/pb/taskpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// maximum 1 megabyte
const maxImageSize = 1 << 20

//TaskServer is the server that provides task services
type TaskServer struct {
	taskpb.UnimplementedTaskServiceServer // implement an interface with private method adding anonymous field for private method
	createTaskUseCase                     contracts.CreateTaskUseCase
	findTaskUseCase                       contracts.FindTaskUseCase
	listTaskUseCase                       contracts.ListTaskUseCase
	loginUseCase                          contracts.LoginUseCase
	imageStore                            ImageStore
}

//NewTaskServer returns a new TaskServer
func NewTaskServer(createTaskUseCase contracts.CreateTaskUseCase,
	findTaskUseCase contracts.FindTaskUseCase,
	listTaskUseCase contracts.ListTaskUseCase,
	loginUseCase contracts.LoginUseCase,
	imageStore ImageStore) *TaskServer {

	return &TaskServer{
		createTaskUseCase: createTaskUseCase,
		findTaskUseCase:   findTaskUseCase,
		listTaskUseCase:   listTaskUseCase,
		loginUseCase:      loginUseCase,
		imageStore:        imageStore,
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
func (t *TaskServer) ListTasks(req *taskpb.ListTaskRequest, stream taskpb.TaskService_ListTasksServer) error {

	err := t.listTaskUseCase.Execute(stream.Context(), req.UserId,
		func(task *entities.Task) error {
			res := &taskpb.ListTaskResponse{
				Task: &taskpb.Task{
					Id:          task.ID,
					UserId:      task.UserID,
					PerformedAt: timestamppb.New(task.PerformedAt),
					Summary:     task.Summary,
				},
			}
			err := stream.Send(res)
			if err != nil {
				return err
			}
			log.Printf("sent task with id: %s", res.Task.Id)

			return nil
		})

	if err != nil {
		return status.Errorf(codes.Internal, "error listing tasks at database: %v", err)
	}

	return nil
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

// UploadImage is client-streaming RPC to upload a task image
func (t *TaskServer) UploadImage(stream taskpb.TaskService_UploadImageServer) error {
	req, err := stream.Recv()
	if err != nil {
		fmt.Printf("cannot receive image info: %v", err)
		return status.Errorf(codes.Unknown, "cannot receive image info")
	}

	taskID := req.GetInfo().GetTaskId()
	imageType := req.GetInfo().GetImageType()

	fmt.Printf("receive an upload-image request for task %s with image type %s", taskID, imageType)

	findTask, err := t.findTaskUseCase.Execute(context.Background(), req.UserId, taskID)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot find task: %v", err)
	}

	if findTask == nil || findTask.Task.ID == "" {
		return status.Errorf(codes.InvalidArgument, "task %s doesn't exist", taskID)
	}

	imageData := bytes.Buffer{}
	imageSize := 0

	for {
		if err := contextError(stream.Context()); err != nil {
			return err
		}

		fmt.Println("waiting to receive more data")

		req, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("no more data")
			break
		}

		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chuck data: %v", err)
		}

		chunk := req.GetChunckData()
		size := len(chunk)

		fmt.Printf("received a chunk with size: %d", size)

		imageSize += size
		if imageSize > maxImageSize {
			return status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize)
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chuck data: %v", err)
		}
	}

	imageID, err := t.imageStore.Save(taskID, imageType, imageData)
	if err != nil {
		return status.Errorf(codes.Internal, "cannot save image to the store: %v", err)
	}

	res := &taskpb.UploadResponse{
		Id:   imageID,
		Size: uint32(imageSize),
	}

	err = stream.SendAndClose(res)
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	fmt.Printf("saved image with id: %s, size: %d\n", imageID, imageSize)

	return nil
}

func contextError(ctx context.Context) error {

	switch ctx.Err() {
	case context.Canceled:
		log.Print("request is canceled")
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		log.Print("deadline is exceeded")
		return status.Error(codes.DeadlineExceeded, "deadline is")
	default:
		return nil
	}

}
