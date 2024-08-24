package main

import (
	"firens/apps/config"
	"firens/apps/controller"
	"firens/apps/pkg/tokens"
	"firens/apps/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := config.ConnectDb()
	if err != nil {
		panic(err)
	}

	r := gin.New()

	r.Use(gin.Logger())

	authController := controller.AuthContorller{
		DB: db,
	}

	v1 := r.Group("/v1")

	auth := v1.Group("auth")

	{
		auth.POST("register", authController.Register)
		auth.POST("login", authController.Login)
		auth.GET("profile", CeckAuth(), authController.Profile)
	}

	r.Run(":4444")
}

func CeckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		bearerToken := strings.Split(header, "Bearer ")

		if len(bearerToken) != 2 {
			resp := response.ResponseAPI{
				StatusCode: http.StatusUnauthorized,
				Message:    "UNAUTHORIZED",
			}

			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		payload, err := tokens.ValidateToken(bearerToken[1])

		if err != nil {
			resp := response.ResponseAPI{
				StatusCode: http.StatusUnauthorized,
				Message:    "Invalid token",
				Payload:    err.Error(),
			}
			ctx.AbortWithStatusJSON(resp.StatusCode, resp)
			return
		}

		ctx.Set("authId", payload.AuthId)

		ctx.Next()
	}
}
