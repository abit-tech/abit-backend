package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func Signup(ctx *gin.Context) {
	// extract the email/pass from body
	var body struct {
		Email    string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// create the user
	user := models.User{
		Email:    body.Email,
		Password: string(hash),
	}
	result := initializers.DB.Create(&user)
	// todo change the error type in case of duplicate email to client error instead of ISE
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	//return
	ctx.JSON(http.StatusOK, gin.H{})
}

func Login(ctx *gin.Context) {
	// extract email and password from body
	var body struct {
		Email    string
		Password string
	}

	if ctx.Bind(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read request body",
		})
		return
	}

	// lookup user from table
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	// if user.ID == 0 {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "invalid email or password",
	// 	})
	// 	return
	// }

	// compare input pass with hash
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	// generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	// return token
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{})
}

func Validate(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"msg": user,
	})
}
