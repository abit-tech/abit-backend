package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"www.github.com/abit-tech/abit-backend/common"
	"www.github.com/abit-tech/abit-backend/initializers"
	"www.github.com/abit-tech/abit-backend/models"
	"www.github.com/abit-tech/abit-backend/utils"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie(common.CookieName)

		authorizationHeader := ctx.Request.Header.Get(common.HeaderKeyAuthorization)
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "not logged in",
			})
			return
		}

		config := initializers.AppConf
		sub, err := utils.ValidateToken(token, config.JWTTokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "bad token",
			})
			return
		}

		var user models.User
		result := initializers.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  "fail",
				"message": "user no longer exists",
			})
			return
		}

		// todo use const
		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
