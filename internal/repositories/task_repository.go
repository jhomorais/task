package repositories

import (
	"context"
	"errors"

	"github.com/fgmaia/task/internal/domain/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) CreateTask(ctx context.Context, entity *entities.Task) error {
	return t.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Create(entity).
		Error
}

func (c *taskRepository) FindTask(ctx context.Context, id string) (*entities.Task, error) {
	var entity *entities.Task

	err := c.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (c *taskRepository) ListTask(ctx context.Context, found func(task *entities.Task) error) error {
	//TODO impl pagination
	var entities []*entities.Task

	err := c.db.
		Preload(clause.Associations).
		Limit(100).
		Order("performed_at desc").
		Find(&entities).Error

	if err != nil {
		return err
	}

	for _, task := range entities {

		if ctx.Err() == context.Canceled {
			return errors.New("request canceled")
		}

		if ctx.Err() == context.DeadlineExceeded {
			return errors.New("deadline is exceeded")
		}

		if err := found(task); err != nil {
			return err
		}
	}

	return nil
}
