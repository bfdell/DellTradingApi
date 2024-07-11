package controllers

import (
	"DellTradingApi/dtos"
	"DellTradingApi/infra"
	"DellTradingApi/models"
	"DellTradingApi/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo: Change query string to fit in requst body in APPEND TICKER
func InitWatchlistRoutes(router *gin.RouterGroup) {
	router.GET("", GetWatchlist)
	router.POST("", AppendTicker)
	router.DELETE("", DeleteTicker)
	router.DELETE("/clear", ClearWatchlist)
}

func GetWatchlist(c *gin.Context) {
	user, err := services.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	//preload watchlist so it can be accessed
	if err := infra.GetDB().Model(&models.UserEntity{}).Preload("Watchlist").First(user).Error; err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	var tickers []string
	for _, listItem := range user.Watchlist {
		tickers = append(tickers, listItem.Ticker)
	}
	fmt.Println(tickers)
	c.JSON(http.StatusOK, gin.H{"tickers": tickers})
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

func DeleteTicker(c *gin.Context) {

}

func ClearWatchlist(c *gin.Context) {

}
