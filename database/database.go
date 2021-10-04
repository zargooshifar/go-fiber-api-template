package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"msgv2-back/config"
	"msgv2-back/models"
	"strconv"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic("DB Connection Failed!")
	}
	fmt.Println("DB Connected!")

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Claims{})
	DB.AutoMigrate(&models.VerificationSMS{})
	DB.AutoMigrate(&models.Contact{})
	fmt.Println("Database Migrated!")
}
