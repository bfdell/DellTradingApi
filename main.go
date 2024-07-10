package main

import (
	"DellTradingApi/controllers"
	"DellTradingApi/infra"
	"DellTradingApi/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//init environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Print("dotenv error", err)
	}

	//init database
	database := infra.OpenDbConnection()
	defer infra.CloseDB(database)
	infra.Migrate(*database)

	ginApi := gin.Default()

	//set up middlewares
	ginApi.Use(middleware.EnsureAuthenticated)

	//init api routes
	routeGroup := ginApi.Group("api/v0")
	controllers.InitUserRoutes(routeGroup.Group("/users"))

	ginApi.Run(":8080")
}
