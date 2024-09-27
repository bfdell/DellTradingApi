package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitPortfolioRoutes(router *gin.RouterGroup) {
	router.GET("", GetPortfolio)
	router.GET("/graph", GetPortfolioGraph)
	router.POST("/buy", BuyStock)
	router.POST("/sell", SellStock)
}

func BuyStock(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var json dtos.PortfolioEntryDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactionErr := services.BuyStock(json.Ticker, json.Shares, user)
	if transactionErr != nil {
		c.JSON(http.StatusBadRequest, transactionErr.Error())
		return
	}

	c.JSON(http.StatusCreated, "")
}

func SellStock(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var json dtos.PortfolioEntryDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactionErr := services.SellStock(json.Ticker, json.Shares, user)
	if transactionErr != nil {
		c.JSON(http.StatusBadRequest, transactionErr.Error())
		return
	}

	c.JSON(http.StatusCreated, "")
}

func GetPortfolio(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	portfolio, portfolioErr := services.GetPortfolioQuotes(user.ID)
	if portfolioErr != nil {
		errors := ""
		for _, error := range portfolioErr {
			errors += error.Error() + "\n"
		}
		c.JSON(http.StatusBadRequest, errors)
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func GetPortfolioGraph(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var timeRange string = c.Query("range")
	if timeRange == "" {
		c.JSON(http.StatusBadRequest, "invalid range")
		return
	}

	graph, graphErr := services.GetPortfolioGraph(user.ID, timeRange)
	if graphErr != nil {
		c.JSON(http.StatusBadRequest, graphErr.Error())
		return
	}

	// fmt.Println("printing graphs")
	// for _, g := range graph {
	// 	fmt.Printf("%+v\t", g)
	// }

	c.JSON(http.StatusOK, graph)
}
