package main

import (
	"firens/apps/config"
	"firens/apps/controller"

	"github.com/gin-gonic/gin"
)

func main(){

	db, err := config.ConnectDb()
	if err != nil{
		panic(err)
	}

	r := gin.New()

	r.Use(gin.Logger())

	authController := controller.AuthContorller{
		DB: db,
	}

	r.POST("/v1/auth/register", authController.Register)
	r.POST("/v1/auth/login", authController.Login)


	r.Run(":4444")
}
