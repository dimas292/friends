package main

import (
	"firens/apps/config"
	"firens/apps/controller"

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
	}

	r.Run(":4444")
}
