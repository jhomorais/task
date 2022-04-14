package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fgmaia/task/pb/taskpb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	PORT = ":8686"
)

func main() {
	conn, err := grpc.Dial("localhost"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	serviceClient := taskpb.NewTaskServiceClient(conn)

	loginRequest := &taskpb.LoginRequest{
		Email:    "technician@gmail.com",
		Password: "tech",
	}
	loginResp, err := serviceClient.Login(context.Background(), loginRequest)
	if err != nil {
		log.Fatal(err)
	}

	if loginResp.UserId == "" {
		log.Fatal("falha ao fazer login")
	}

	summary := "created clean arch project"

	argsWithoutProg := os.Args[1:] //argsWithProg := os.Args

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] != "" {
		summary = argsWithoutProg[0]
	}

	//testCreateTask(serviceClient, loginResp.UserId, summary)
	testUploadImage(serviceClient, loginResp.UserId, summary)
}

func uploadImage(serviceClient taskpb.TaskServiceClient, userId string, taskId string, imagePath string) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stream, err := serviceClient.UploadImage(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	req := &taskpb.UploadImageRequest{
		UserId: userId,
		Data: &taskpb.UploadImageRequest_Info{
			Info: &taskpb.ImageInfo{
				TaskId:    taskId,
				ImageType: filepath.Ext(imagePath),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		log.Fatal("cannot send image info: ", err, stream.RecvMsg(nil))
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("cannot read chunck to buffer: ", err)
		}

		req := &taskpb.UploadImageRequest{
			Data: &taskpb.UploadImageRequest_ChunckData{
				ChunckData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			err2 := stream.RecvMsg(nil)
			log.Fatal("cannot send chunk to server: ", err, err2)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("image uploaded with id: %s, size: %d", res.GetId(), res.GetSize())
}

func testUploadImage(serviceClient taskpb.TaskServiceClient, userId string, summary string) {
	taskId := testCreateTask(serviceClient, userId, summary)
	uploadImage(serviceClient, userId, taskId, "../../how_to_test_console.png")
}

func testListTasks(serviceClient taskpb.TaskServiceClient, userId string) {

}

func testCreateTask(serviceClient taskpb.TaskServiceClient, userId string, summary string) string {

	task := &taskpb.Task{
		Summary:     summary,
		PerformedAt: timestamppb.New(time.Now()),
		UserId:      userId,
	}

	req := &taskpb.CreateTaskRequest{
		Task: task,
	}

	res, err := serviceClient.CreateTask(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task Created ID: " + res.Id)

	return res.Id
}
