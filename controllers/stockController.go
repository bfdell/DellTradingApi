package controllers

import (
	"DellTradingApi/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitStockRoutes(router *gin.RouterGroup) {
	router.GET("", GetQuote)
}

// ! This endpoint is not used
func GetQuote(c *gin.Context) {
	_, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	var ticker string = c.Query("ticker")
	if ticker == "" {
		c.JSON(http.StatusBadRequest, "invalid ticker")
		return
	}

	if quote, quoteErr := services.GetQuote(ticker); quoteErr != nil {
		c.JSON(http.StatusUnprocessableEntity, quoteErr.Error())
		return
	} else {
		c.JSON(http.StatusOK, quote)
	}
}
