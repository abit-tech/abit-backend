package initializers

import (
	"fmt"

	"www.github.com/abit-tech/abit-backend/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Perk{})
	DB.AutoMigrate(&models.Video{})
	DB.AutoMigrate(&models.Token{})
	DB.AutoMigrate(&models.Waitlist{})
	fmt.Println("🚀 Database migration complete")
}
