package di

import (
	"fmt"

	"github.com/fgmaia/task/config"
	"github.com/fgmaia/task/internal/domain/entities"
	"github.com/fgmaia/task/sample"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitGormMysqlDB() (*gorm.DB, error) {
	config.LoadServerEnvironmentVars()

	dsn := fmt.Sprintf("%s:%s@%s", config.GetMysqlUser(), config.GetMysqlPassword(), config.GetMysqlConnectionString())

	mysqlDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	mysqlDb.AutoMigrate(&entities.Task{}, &entities.User{})

	sample.DBSeed(mysqlDb)

	return mysqlDb, err
}
