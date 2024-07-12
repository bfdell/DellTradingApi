package services

import (
	"DellTradingApi/infra"
	"DellTradingApi/models"
	"fmt"

	"gorm.io/gorm"
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
	}

	db := infra.GetDB()
	result := db.Create(newWatchlistItem)

	if result.Error != nil {
		//gets already existing ticker if exists
		result := db.Unscoped().Model(&models.WatchlistEntity{}).Find(newWatchlistItem)
		if result.Error != nil {
			return result.Error
		}

		//if existing ticker has been deleted before, unmark it as deleted
		if newWatchlistItem.DeletedAt.Valid {
			newWatchlistItem.DeletedAt = gorm.DeletedAt{}
			err := db.Save(newWatchlistItem).Error
			return err
		} else {
			return fmt.Errorf("ticker already exists in watchlist")
		}
	}

	return result.Error
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

func ClearWatchListItems(user *models.UserEntity) error {
	result := infra.GetDB().Where("user_id = ?", user.ID).Delete(&models.WatchlistEntity{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("watchlist was empty")
	}

	return nil
}
