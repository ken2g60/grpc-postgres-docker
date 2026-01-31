package database

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func InitDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	dsn := "postgresql://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}
	return db

}

type Migrations struct {
	DB     *gorm.DB
	Models []interface{}
}

func RunMigrations(migrations Migrations) {
	for _, model := range migrations.Models {
		err := migrations.DB.AutoMigrate(model)
		if err != nil {
			logrus.Error("Error in migration", zap.Error(err))

		}
	}
}
