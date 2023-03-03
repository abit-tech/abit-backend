package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"www.github.com/abit-tech/abit-backend/controllers"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/middleware"
)

var server *gin.Engine

func init() {
	err := initializers.LoadConfig(".")
	if err != nil {
		panic("error in reading config")
	}

	initializers.ConnectDB()
	initializers.SyncDatabase()
	initializers.SetupGoogleOauth()

	server = gin.Default()

	prettyConf, _ := json.MarshalIndent(initializers.AppConf, "", "\t")
	fmt.Printf("config object: %v\n", string(prettyConf))
}

func main() {
	r := gin.Default()
	r.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "server is up and running"})
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
	}

	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "server is up and running",
		})
	})

	auth_router := router.Group("/auth")
	auth_router.POST("/register", controllers.SignUpUser)
	auth_router.POST("/login", controllers.SignInUser)
	auth_router.GET("/logout", middleware.DeserializeUser(), controllers.LogoutUser)

	google_router := router.Group("/google")
	google_router.GET("/login", controllers.GoogleLogin)

	router.GET("/sessions/oauth/google", controllers.GoogleOAuth)
	router.GET("/users/me", middleware.DeserializeUser(), controllers.GetMe)

	router.StaticFS("/images", http.Dir("public"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "route not found",
		})
	})

	// r.POST("/signup", controllers.Signup)
	// r.POST("/login", controllers.Login)
	// r.POST("/validate", middleware.RequireAuth, controllers.Validate)

	// auth_router := r.Group("/google")
	// auth_router.GET("/register", controllers.GoogleLogin)
	// auth_router.GET("/callback", controllers.GoogleCallback)
	// // auth_router.POST("/login", controllers.GoogleSignInUser)
	// // auth_router.GET("/logout", middleware.DeserializeUser(), controllers.GoogleLogoutUser)

	server.Run(":" + "3000")
}
