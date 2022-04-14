package main

import (
	"fmt"
	"log"
	"net"

	"github.com/fgmaia/task/internal/infra/di"
	"github.com/fgmaia/task/internal/services"
	"github.com/fgmaia/task/pb/taskpb"
	"google.golang.org/grpc"
)

const (
	PORT = ":8686"
)

func main() {
	dependencies := di.NewBuild()

	imageStore := services.NewDiskImageStore("img")

	taskServer := services.NewTaskServer(dependencies.Usecases.CreateTaskUseCase,
		dependencies.Usecases.FindTaskUseCase,
		dependencies.Usecases.ListTaskUseCase,
		dependencies.Usecases.LoginUseCase,
		imageStore)

	grpcServer := grpc.NewServer()
	taskpb.RegisterTaskServiceServer(grpcServer, taskServer)

	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GRPC SERVER LISTEN PORT: " + PORT)
	grpcServer.Serve(listener)
}
