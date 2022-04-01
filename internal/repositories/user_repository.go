package repositories

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (t *userRepository) CreateUser(ctx context.Context, entity *entities.User) error {
	entity.Password = utils.HashPassword(entity.Password, []byte(utils.SALT))
	return t.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Save(&entity).
		Error
}

func (c *userRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	var entity *entities.User

	err := c.db.
		Preload(clause.Associations).
		Where("id = ?", id).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (c *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var entity *entities.User

	err := c.db.
		Preload(clause.Associations).
		Where("email = ?", email).
		Limit(1).
		Find(&entity).Error

	return entity, err
}

func (c *userRepository) UpdateUser(ctx context.Context, enitity *entities.User) error {
	err := c.db.
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&enitity).
		Error

	return err
}
