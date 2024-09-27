package services

import (
	"DellTradingApi/dtos"
	"DellTradingApi/infra"
	"DellTradingApi/models"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func FindLastTransaction(ticker string, ID uint) (*models.PortfolioEntity, error) {
	newPortfolioEntry := &models.PortfolioEntity{
		UserID: ID,
		Ticker: ticker,
	}

	//get most recent portfolio entry of input ticker from user
	queryResult := infra.GetDB().Model(&models.PortfolioEntity{}).Order("created_at desc").First(newPortfolioEntry).Error

	//adjust createdAt variable just in case this entry is old to create a new one
	newPortfolioEntry.CreatedAt = time.Now()
	return newPortfolioEntry, queryResult
}

// todo: what if the user waits until the stock changes to buy?!
func BuyStock(ticker string, shares uint, user *models.UserEntity) error {
	portfolioEntry, err := FindLastTransaction(ticker, user.ID)
	//if we could not find a previous record, do not return becuase we are buying this stock for the first time
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	stockResponse, err := GetQuote(ticker)
	if err != nil {
		return err
	}

	//make sure I have enough money to buy that much
	var stockPrice float64 = float64(stockResponse.Price)
	cost := stockPrice * float64(shares)
	if cost > user.Cash {
		return fmt.Errorf("not enough money")
	}

	user.Cash -= cost
	portfolioEntry.Shares += shares

	userErr := infra.GetDB().Save(user).Error
	if userErr != nil {
		return userErr
	}
	portfolioErr := infra.GetDB().Save(portfolioEntry).Error
	if portfolioErr != nil {
		return portfolioErr
	}

	return nil
}

func SellStock(ticker string, shares uint, user *models.UserEntity) error {
	portfolioEntry, err := FindLastTransaction(ticker, user.ID)
	//if a former entry is not found in this case, we return because there must be shares of a stock to sell it
	if err != nil {
		fmt.Println("THE ERROR", err.Error())
		return err
	}

	stockResponse, err := GetQuote(ticker)
	if err != nil {
		return err
	}

	//make sure I have enough shares to sell that much
	var stockPrice float64 = float64(stockResponse.Price)
	if shares > portfolioEntry.Shares {
		return fmt.Errorf("not enough shares")
	}

	revenue := stockPrice * float64(shares)
	user.Cash += revenue
	portfolioEntry.Shares -= shares

	userErr := infra.GetDB().Save(user).Error
	if userErr != nil {
		return userErr
	}
	portfolioErr := infra.GetDB().Save(portfolioEntry).Error
	if portfolioErr != nil {
		return portfolioErr
	}

	return nil
}

func GetPortfolio(ID uint) ([]*dtos.PortfolioEntryDto, error) {
	//returns every ticker owned and how many shares
	query := `
	SELECT ticker, shares from portfolio_entities 
	INNER JOIN (SELECT MAX(created_at) as most_recent from portfolio_entities GROUP BY ticker) 
	as portfolio2 on created_at = portfolio2.most_recent and user_id = ?
	`

	var results []*dtos.PortfolioEntryDto
	err := infra.GetDB().Raw(query, ID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetPortfolioQuotes(ID uint) ([]*dtos.PortfolioAssetDto, []error) {
	portfolioEntries, err := GetPortfolio(ID)
	if err != nil {
		return nil, []error{err}
	}

	var portfolioAssets []*dtos.PortfolioAssetDto
	var quoteErrors []error
	for _, entry := range portfolioEntries {
		if entry.Shares > 0 {
			quote, quoteErr := GetQuote(entry.Ticker)

			//dont inlude the stocks that have trouble fetching
			if quoteErr != nil {
				quoteErrors = append(quoteErrors, quoteErr)
			} else {
				asset := &dtos.PortfolioAssetDto{
					StockQuoteDto: *quote,
					Shares:        entry.Shares,
				}
				portfolioAssets = append(portfolioAssets, asset)
			}
		}
	}

	if len(quoteErrors) == 0 {
		quoteErrors = nil
	}

	return portfolioAssets, quoteErrors
}
