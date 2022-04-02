package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fgmaia/task/config"
	"github.com/fgmaia/task/internal/infra/queue"
)

func main() {
	config.LoadServerEnvironmentVars()

	chQuit := make(chan os.Signal, 2)
	signal.Notify(chQuit, os.Interrupt, syscall.SIGTERM)

	go func() {
		for range chQuit {
			os.Exit(0)
		}
	}()

	taskQueue := queue.NewRabbitMQ(config.GetRabbitMQClient(),
		"task-exchange",
		"tasks")

	err := taskQueue.InitQueue(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for {
			wg.Add(1)
			go func() {
				defer wg.Done()
				consumeResponseQueue(context.Background(), taskQueue)
			}()
			wg.Wait()
			time.Sleep(time.Second * 5)
		}
	}()

	<-forever
}

func consumeResponseQueue(ctx context.Context, taskQueue queue.RabbitMQ) {
	messages, err := taskQueue.Consume(ctx)
	if err != nil {
		err = taskQueue.InitQueue(ctx)
		if err != nil {
			log.Fatal(err)
		}

		messages, err = taskQueue.Consume(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("WAITING TO RECEIVE MESSAGES")
	for message := range messages {
		fmt.Printf("RECEIVED MESSAGE: %s \n", string(message.Body))
	}
}
