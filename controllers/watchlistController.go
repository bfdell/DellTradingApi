package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitWatchlistRoutes(router *gin.RouterGroup) {
	router.GET("", GetWatchlist)
	router.POST("/append", AppendTicker)
	router.DELETE("/remove", RemoveTicker)
	router.DELETE("/clear", ClearWatchlist)
}

func GetWatchlist(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	if tickers, loadErr := services.GetWatchlistItems(user); loadErr != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	} else {
		fmt.Println(tickers)
		c.JSON(http.StatusOK, gin.H{"tickers": tickers})
	}
}

func AppendTicker(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve ticker from request body
	var json dtos.WatchlistRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		//todo: more verbose json errors
		c.JSON(http.StatusBadRequest, err)
		return
	}

	createErr := services.CreateWatchlistItem(json.Ticker, user.ID)
	if createErr != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	c.JSON(http.StatusCreated, "")
}

func RemoveTicker(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Retrieve ticker from request body
	var json dtos.WatchlistRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		//todo: more verbose json errors
		c.JSON(http.StatusBadRequest, err)
		return
	}

	removeErr := services.RemoveWatchlistItem(json.Ticker, user.ID)
	if removeErr != nil {
		c.JSON(http.StatusUnprocessableEntity, removeErr.Error())
		return
	}

	c.JSON(http.StatusOK, "")
}

func ClearWatchlist(c *gin.Context) {

}
