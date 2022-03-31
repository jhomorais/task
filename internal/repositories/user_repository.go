package repositories

import (
	"context"

	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/internal/repositories/contracts"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contracts.UserRepository {
	return &userRepository{db: db}
}

func (t *userRepository) CreateUser(ctx context.Context, entity *entities.User) error {
	return t.db.
		Session(&gorm.Session{FullSaveAssociations: false}).
		Save(entity).
		Error
}

/*

func (c *chargeRepository) CreateCharge(chargeModel *models.Charge) error {
	err := c.db.Client().
		Create(&chargeModel).
		Error

	return err
}

func (c *chargeRepository) UpdateCharge(chargeModel *models.Charge) error {
	err := c.db.Client().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&chargeModel).
		Error

	return err
}

func (c *chargeRepository) FindChargeByReferenceKey(referenceKey string) (*models.Charge, error) {
	var chargeModel *models.Charge

	err := c.db.Client().
		Preload(clause.Associations).
		Last(&chargeModel, "reference_key = ?", referenceKey).
		Error

	return chargeModel, err
}

func (c *chargeRepository) FindChargeById(id uint) (*models.Charge, error) {
	var chargeModel *models.Charge

	err := c.db.Client().
		Preload(clause.Associations).
		Last(&chargeModel, "id = ?", id).
		Error

	return chargeModel, err
}

func (c *chargeRepository) FindChargeByExternalId(externalId string) (*models.Charge, error) {
	var chargeModel *models.Charge

	err := c.db.Client().
		Preload(clause.Associations).
		Last(&chargeModel, "external_id = ?", externalId).
		Error

	return chargeModel, err
}

func (c *chargeRepository) FindExpireds(expireAt time.Time) ([]models.Charge, error) {
	var charges []models.Charge

	err := c.db.Client().
		Limit(120).
		Where("internal_status = 'waiting_payment' AND pix_expiration_at <= ?", expireAt).
		Order("pix_expiration_at asc").
		Find(&charges).
		Error

	return charges, err
}

*/
