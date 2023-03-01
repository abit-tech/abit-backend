package main

import (
	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/controllers"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}
func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.POST("/validate", middleware.RequireAuth, controllers.Validate)
	r.Run()
}
