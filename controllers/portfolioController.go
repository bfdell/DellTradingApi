package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitPortfolioRoutes(router *gin.RouterGroup) {
	router.POST("/buy", BuyStock)
	router.POST("/sell", SellStock)
}

func BuyStock(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var json dtos.PortfolioUpdateRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		//todo: more verbose json errors
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactionErr := services.BuyStock(json.Ticker, json.Shares, user)
	if transactionErr != nil {
		fmt.Println("RECIEVED TRANSACTION ERR")
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

	var json dtos.PortfolioUpdateRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		//todo: more verbose json errors
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
