package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitPortfolioRoutes(router *gin.RouterGroup) {
	router.GET("", GetPortfolio)
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
