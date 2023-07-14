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
		fmt.Printf("could not find file, reading from ENV")
	}

	prettyConf, _ := json.MarshalIndent(initializers.AppConf, "", "\t")
	fmt.Printf("config object: %v\n", string(prettyConf))
	initializers.ConnectDB()
	initializers.SyncDatabase()
	initializers.SetupGoogleOauth()

	server = gin.Default()

}

func main() {
	r := gin.Default()
	r.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "server is up and running"})
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"https://getabit.co",
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

	router.POST("/waitlist", controllers.AddToWaitlist)

	auth_router := router.Group("/auth")
	auth_router.POST("/register", controllers.SignUpUser)
	auth_router.POST("/login", controllers.SignInUser)
	auth_router.GET("/logout", middleware.DeserializeUser(), controllers.LogoutUser)

	google_router := router.Group("/google")
	google_router.GET("/login/user", controllers.GoogleLoginForUser)
	google_router.GET("/login/creator", controllers.GoogleLoginForCreator)

	router.GET("/sessions/oauth/google", controllers.GoogleOAuth)
	// router.GET("/sessions/oauth/google/creator", controllers.GoogleOAuthForCreator)
	router.GET("/users/me", middleware.DeserializeUser(), controllers.GetMe)

	video_router := router.Group("/video")
	video_router.POST("/create", middleware.DeserializeUser(), controllers.CreateVideo)
	video_router.POST("/fetch", middleware.DeserializeUser(), controllers.GetVideo)
	video_router.POST("/update", middleware.DeserializeUser(), controllers.UpdateVideo)
	video_router.POST("/delete", middleware.DeserializeUser(), controllers.DeleteVideo)
	video_router.POST("/buytoken", middleware.DeserializeUser(), controllers.PurchaseToken)

	dash_router := router.Group("/dash")
	dash_router.POST("/user", middleware.DeserializeUser(), controllers.UserDash)
	dash_router.POST("/creator", middleware.DeserializeUser(), controllers.CreatorDash)

	// creator_router := router.Group("/creator")
	// creator_router.POST("/profile", controllers.CreatorDummy)
	// creator_router.POST("/videos", controllers.CreatorDummy)

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

	server.Run()
}
