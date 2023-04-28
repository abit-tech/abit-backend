package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
)

func AddToWaitlist(ctx *gin.Context) {
	var payload *models.WaitlistInput
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "bad request argument",
		})
		return
	}

	newEntry := models.Waitlist{
		Email: payload.Email,
	}

	result := initializers.DB.Create(&newEntry)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": "something went wrong",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "thank you for your interest!",
	})
}
