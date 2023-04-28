package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"www.github.com/abit-tech/abit-backend/common"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
	"www.github.com/abit-tech/abit-backend/utils"
)

// todo should use constants for status values in response

func SignUpUser(ctx *gin.Context) {
	var payload *models.RegisterUserInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// check for valid role
	if payload.Role != common.RoleUser &&
		payload.Role != common.RoleAdmin &&
		payload.Role != common.RoleCreator {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "invalid role",
		})
		return
	}

	// if creator registration, ensure channel link is passed
	if payload.Role == common.RoleCreator && payload.ChannelLink == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "provide channel link for creator sign up",
		})
		return
	}

	newUser := models.User{
		Name:        payload.Name,
		Email:       strings.ToLower(payload.Email),
		Password:    string(hash),
		Provider:    common.ProviderLocal,
		ChannelLink: payload.ChannelLink,
		Role:        payload.Role,
		Verified:    payload.Role == common.RoleUser, // true for users by default, false for all others
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

	if user.Provider == common.ProviderGoogle {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "use google login instead",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "invalid email or password",
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
	config := initializers.AppConf
	ctx.SetCookie(common.CookieName, "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "logged out",
	})
	ctx.Redirect(http.StatusPermanentRedirect, config.FrontEndOrigin)
}

func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	// var pathURL string = "/"

	// if ctx.Query("state") != "" {
	// 	pathURL = ctx.Query("state")
	// }

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

	role, err := ctx.Cookie(common.RoleCookie)
	if err != nil {
		fmt.Print("role not found in google login context")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "something went wrong",
		})
		return
	}

	// unset the role cookie
	ctx.SetCookie(common.RoleCookie, "", -1, "/", "localhost", false, true)

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	userData := models.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		Photo:     googleUser.Picture,
		Provider:  common.ProviderGoogle,
		Role:      role,
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
	// ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathURL))
	ctx.Redirect(http.StatusTemporaryRedirect, "www.facebook.com")
}
