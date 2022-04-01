package di

import (
	"context"
	"log"

	"github.com/fgmaia/task/config"
	"github.com/fgmaia/task/internal/infra/queue"
	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"gorm.io/gorm"
)

type DenpencyBuild struct {
	DB           *gorm.DB
	TaskQueue    queue.RabbitMQ
	Repositories Repositories
	Usecases     Usecases
}

type Repositories struct {
	UserRepository repositories.UserRepository
	TaskRepository repositories.TaskRepository
}

type Usecases struct {
	CreateTaskUseCase contracts.CreateTaskUseCase
	FindTaskUseCase   contracts.FindTaskUseCase
	ListTaskUseCase   contracts.ListTaskUseCase
}

func NewBuild() *DenpencyBuild {

	builder := &DenpencyBuild{}

	builder = builder.buildDB().
		buildQueue().
		buildRepositories().
		buildUseCases()

	return builder
}

func (d *DenpencyBuild) buildDB() *DenpencyBuild {
	var err error
	d.DB, err = InitGormMysqlDB()
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func (d *DenpencyBuild) buildQueue() *DenpencyBuild {
	d.TaskQueue = queue.NewRabbitMQ(config.GetRabbitMQClient(),
		"task-exchange",
		"tasks")

	err := d.TaskQueue.InitQueue(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func (d *DenpencyBuild) buildRepositories() *DenpencyBuild {
	d.Repositories.UserRepository = repositories.NewUserRepository(d.DB)
	d.Repositories.TaskRepository = repositories.NewTaskRepository(d.DB)
	return d
}

func (d *DenpencyBuild) buildUseCases() *DenpencyBuild {
	d.Usecases.CreateTaskUseCase = usecases.NewCreateTaskUseCase(d.Repositories.UserRepository,
		d.Repositories.TaskRepository,
		d.TaskQueue)
	d.Usecases.FindTaskUseCase = usecases.NewFindTaskUseCase(d.Repositories.UserRepository, d.Repositories.TaskRepository)
	d.Usecases.ListTaskUseCase = usecases.NewListTaskUseCase(d.Repositories.UserRepository, d.Repositories.TaskRepository)

	return d
}
