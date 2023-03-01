package initializers

import "www.github.com/abit-tech/abit-backend/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
