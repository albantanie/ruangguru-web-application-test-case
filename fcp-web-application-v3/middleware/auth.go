package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")

		if ctx.Request.URL.Path == "/user/login" || ctx.Request.URL.Path == "/user/register" {
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
			if err == http.ErrNoCookie {
				if ctx.Request.Header.Get("Content-Type") == "application/json" {
					ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
					return
				}
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorResponse{Error: "Bad Request"})
			ctx.Abort()
			return
		}

		tokenClaim := model.Claims{}
		token, err := jwt.ParseWithClaims(cookie.Value, &tokenClaim, func(t *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})

		if err != nil {
			if ctx.Request.Header.Get("Content-Type") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
				return
			}
			ctx.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Bad Request"})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Unauthorized"})
			return
		}

		ctx.Set("email", tokenClaim.Email)
		// TODO: answer here
	})
}
