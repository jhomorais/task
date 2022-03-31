package repositories

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories/contracts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) contracts.TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) CreateTask(ctx context.Context, entity *entities.Task) error {
	return t.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Save(entity).
		Error
}

func (c *taskRepository) FindTask(ctx context.Context, id string) (*entities.Task, error) {
	var entity *entities.Task

	err := c.db.
		Preload(clause.Associations).
		Last(&entity, "id = ?", id).
		Error

	return entity, err
}
