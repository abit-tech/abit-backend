package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/common"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
	"www.github.com/abit-tech/abit-backend/utils"
)

// todo should use constants for status values in response

func SignUpUser(ctx *gin.Context) {
	var payload *models.User
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	now := time.Now()
	newUser := models.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		Role:      "user", // todo this might change as we introduce creators
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := initializers.DB.Create(&newUser)

	if result.Error != nil && strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
		ctx.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "user with that email already exists",
		})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user": models.FilteredResponse(&newUser),
		},
	})
}

func SignInUser(ctx *gin.Context) {
	var payload *models.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid email or password",
		})
		return
	}

	if user.Provider == "Google" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "use google login instead",
		})
		return
	}

	config := initializers.AppConf
	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	ctx.SetCookie(common.CookieName, token, config.TokenMaxAge*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged in",
	})
}

func LogoutUser(ctx *gin.Context) {
	ctx.SetCookie(common.CookieName, "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged out",
	})
}

func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathURL string = "/"

	if ctx.Query("state") != "" {
		pathURL = ctx.Query("state")
	}

	if code == "" {
		// todo log error
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "authorization code not provided",
		})
		return
	}

	tokenRes, err := utils.GetGoogleOauthToken(code)
	if err != nil {
		// todo log error
		fmt.Printf("error in getting token: %v\n", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	googleUser, err := utils.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		// todo log error
		fmt.Printf("error in getting user: %v\n", err.Error())
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	userData := models.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		Photo:     googleUser.Picture,
		Provider:  "Google", // todo use const
		Role:      "user",   // todo decide this
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if initializers.DB.Model(&userData).Where("email = ?", email).Updates(&userData).RowsAffected == 0 {
		initializers.DB.Create(&userData)
	}

	var user models.User
	initializers.DB.First(&user, "email = ?", email)

	config := initializers.AppConf
	token, err := utils.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		// todo log error
		fmt.Printf("error in getting jwt token: %v\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	ctx.SetCookie(common.CookieName, token, config.TokenMaxAge*60, "/", "localhost", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathURL))
}
