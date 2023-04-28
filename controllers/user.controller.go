package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/models"
)

func GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": models.FilteredResponse(&currentUser),
		},
	})
}
