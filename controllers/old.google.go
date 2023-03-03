package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
	"www.github.com/abit-tech/abit-backend/utils"
)

/*
The same-origin-allow-popups value is recommended for the Cross-Origin-Opener-Policy header on pages where Sign In With Google button and/or Google One Tap are displayed.

Set the Referrer-Policy: no-referrer-when-downgrade header when testing using http and localhost.

*/

func GoogleLogin(ctx *gin.Context) {
	url := initializers.GoogleConfig.AuthCodeURL("randomstate")
	ctx.Redirect(http.StatusTemporaryRedirect, url)

}

func GoogleCallback(ctx *gin.Context) {
	// code := ctx.Query("code")
	// var pathUrl string = "/"

	// if ctx.Query("state") != "" {
	// 	pathUrl = ctx.Query("state")
	// }

	// if code == "" {
	// 	ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
	// 	return
	// }

	state := ctx.Query("state")
	if state != "randomstate" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "state mismatch",
		})
		return
	}

	code := ctx.Query("code")
	tokenRes, err := utils.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	google_user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fmt.Printf("retrieved google user: %v\n", google_user)
}

/// check if this is good or should be deleted?

func GoogleSignUpUser(ctx *gin.Context) {
	var payload *models.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// now := time.Now()
	newUser := models.User{
		// Name:      payload.Name,
		Email:    strings.ToLower(payload.Email),
		Password: payload.Password,
		// Role:      "user",
		// Verified:  true,
		// CreatedAt: now,
		// UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": newUser}})
}

func GoogleSignInUser(ctx *gin.Context) {
	var payload *models.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	tokenSecret := os.Getenv("JWT_SECRET")

	token, err := utils.GenerateToken(1*time.Hour, user.ID, tokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, 3660, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GoogleLogoutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func GoogleOAuthOld(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	tokenRes, err := utils.GetGoogleOauthToken(code)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	google_user, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// now := time.Now()
	email := strings.ToLower(google_user.Email)

	user_data := models.User{
		// Name:      google_user.Name,
		Email:    email,
		Password: "",
		// Photo:     google_user.Picture,
		// Provider:  "Google",
		// Role:      "user",
		// Verified:  true,
		// CreatedAt: now,
		// UpdatedAt: now,
	}

	if initializers.DB.Model(&user_data).Where("email = ?", email).Updates(&user_data).RowsAffected == 0 {
		initializers.DB.Create(&user_data)
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	tokenSecret := os.Getenv("JWT_SECRET")

	token, err := utils.GenerateToken(1*time.Hour, user.ID, tokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint("http://localhost:3000", pathUrl))
}
