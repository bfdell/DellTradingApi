package main

import (
	"DellTradingApi/controllers"
	"DellTradingApi/infra"
	"DellTradingApi/middleware"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// todo: remove verbose error returning after everything maybe?
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

	//init stock service
	// services.InitStockService()

	//gin routing
	ginApi := gin.Default()

	// Add the CORS middleware
	ginApi.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//set up middlewares
	ginApi.Use(middleware.EnsureAuthenticated)

	//init api routes
	routeGroup := ginApi.Group("api/v0")
	controllers.InitUserRoutes(routeGroup.Group("/users"))
	controllers.InitWatchlistRoutes(routeGroup.Group("/watchlist"))
	controllers.InitPortfolioRoutes(routeGroup.Group("/portfolio"))
	controllers.InitStockRoutes(routeGroup.Group("/stock"))

	ginApi.Run(":8080")
}
