package initializers

import (
	"fmt"

	"www.github.com/abit-tech/abit-backend/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	fmt.Println("🚀 Database migration complete")
}
