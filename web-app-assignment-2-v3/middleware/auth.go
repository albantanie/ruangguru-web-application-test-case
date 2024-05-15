package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")

		if ctx.Request.URL.Path == "/user/register" || ctx.Request.URL.Path == "/user/login" {
			ctx.Next()
			return
		}

		if err != nil {
			if err == http.ErrNoCookie {
				if ctx.Request.Header.Get("Content-Type") == "application/json" {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "Unauthorized",
					})
					return
				}
				ctx.Redirect(http.StatusSeeOther, "/login")
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			ctx.Abort()
			return
		}

		tokenClaims := model.Claims{}
		token, err := jwt.ParseWithClaims(cookie, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "Unauthorized",
				})
				return
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad Request",
			})
			ctx.Abort()
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		ctx.Set("id", tokenClaims.UserID)
		ctx.Next()
	})
}
