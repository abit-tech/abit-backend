package initializers

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		AppConf.DBHost,
		AppConf.DBUser,
		AppConf.DBPassword,
		AppConf.DBName,
		AppConf.DBPort,
		AppConf.DBSSLMode)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("could not connect to DB")
	}
	fmt.Println("🚀 Connected Successfully to the Database")
}
