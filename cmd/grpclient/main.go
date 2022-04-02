package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	task := &taskpb.Task{
		Summary:     summary,
		PerformedAt: timestamppb.New(time.Now()),
		UserId:      loginResp.UserId,
	}

	req := &taskpb.CreateTaskRequest{
		Task: task,
	}

	res, err := serviceClient.CreateTask(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task Created ID: " + res.Id)
}
