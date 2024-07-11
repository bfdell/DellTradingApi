package services

import (
	"DellTradingApi/infra"
	"DellTradingApi/models"
	"fmt"
)

func GetWatchlistItems(user *models.UserEntity) ([]string, error) {
	//preload watchlist so it can be accessed
	err := infra.GetDB().Model(&models.UserEntity{}).Preload("Watchlist").First(user).Error

	var tickers []string
	if err == nil {
		for _, listItem := range user.Watchlist {
			tickers = append(tickers, listItem.Ticker)
		}
	}

	return tickers, err
}

func CreateWatchlistItem(ticker string, ID uint) error {
	newWatchlistItem := &models.WatchlistEntity{
		UserID: ID,
		Ticker: ticker,
		// CreatedAt: time.Now(),
	}

	return infra.GetDB().Create(newWatchlistItem).Error
}

func RemoveWatchlistItem(ticker string, ID uint) error {
	result := infra.GetDB().Model(&models.WatchlistEntity{}).Delete(&models.WatchlistEntity{
		Ticker: ticker, UserID: ID})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("ticker does not exist")
	}

	return nil
}
