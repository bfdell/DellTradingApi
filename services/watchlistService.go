package services

import (
	"DellTradingApi/infra"
	"DellTradingApi/models"
)

func CreateWatchlistItem(ticker string, ID uint) error {
	newWatchlistItem := &models.WatchlistEntity{
		UserID: ID,
		Ticker: ticker,
		// CreatedAt: time.Now(),
	}

	return infra.GetDB().Create(newWatchlistItem).Error
}
