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

func getRawPortfolio(ID uint) (map[string][]*models.PortfolioEntity, []error) {
	portfolioTickers, err := GetPortfolio(ID)
	if err != nil {
		return nil, []error{err}
	}

	db := infra.GetDB()
	assetMap := make(map[string][]*models.PortfolioEntity)
	var assetErrors []error
	for _, tickerEntry := range portfolioTickers {

		var tickers []*models.PortfolioEntity
		queryErr := db.Model(&models.PortfolioEntity{}).Where("ticker = ? AND user_id = ?", tickerEntry.Ticker, ID).Order("created_at").Find(&tickers).Error
		if queryErr != nil {
			assetErrors = append(assetErrors, queryErr)
		} else {
			assetMap[tickerEntry.Ticker] = tickers
			// fmt.Println("found tickers for", tickerEntry.Ticker, " ")
			// for _, t := range tickers {
			// 	fmt.Printf("%+v\t", t)
			// }
			// fmt.Println("\n\n\n")
		}
	}

	if len(assetErrors) == 0 {
		assetErrors = nil
	}

	return assetMap, assetErrors
}

func getStockHistory(portfolio map[string][]*models.PortfolioEntity, startDate string) (map[string]map[string]*dtos.TimeSeriesQuoteDto, []error) {
	tickerHistoryMap := make(map[string]map[string]*dtos.TimeSeriesQuoteDto)
	var tickerErrors []error
	for key := range portfolio {
		tickerHistory, err := GetHistory(key, startDate)
		if err != nil {
			tickerErrors = append(tickerErrors, err)
		} else {
			tickerHistoryMap[key] = tickerHistory
		}
	}

	if len(tickerErrors) == 0 {
		tickerErrors = nil
	}

	return tickerHistoryMap, tickerErrors
}

// todo: add support for day
func GetPortfolioGraph(ID uint, timeRange string) ([]*dtos.PortfolioGraphDto, error) {
	portfolioEntities, portfolioErrs := getRawPortfolio(ID)
	if portfolioErrs != nil {
		// for _, err := range portfolioErrs {
		// 	fmt.Printf("%+v\n", err)
		// }
		return nil, fmt.Errorf("failed to retrive portfolio data")
	}

	var graphData []*dtos.PortfolioGraphDto
	var yearOffset, monthOffset, dayOffset = 0, 0, 0
	switch timeRange {
	case "week":
		dayOffset = -7
	case "month":
		monthOffset = -1
	case "year":
		yearOffset = -1
	default:
		return nil, fmt.Errorf("invalid time range")
	}

	now := time.Now()
	today := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
	startDate := today.AddDate(yearOffset, monthOffset, dayOffset)
	dayOfWeek := startDate.Weekday()
	startDateStr := startDate.Format("2006-01-02")
	if dayOfWeek == time.Saturday {
		startDateStr = startDate.AddDate(0, 0, -1).Format("2006-01-02")
	}
	if dayOfWeek == time.Sunday {
		startDateStr = startDate.AddDate(0, 0, -2).Format("2006-01-02")
	}
	tickerHistory, tickerErrs := getStockHistory(portfolioEntities, startDateStr)
	if tickerErrs != nil {
		// for _, err := range tickerErrs {
		// 	fmt.Printf("%+v\n", err)
		// }
		return nil, fmt.Errorf("failed to retrive ticker history data")
	}

	// fmt.Println("printing ticker history")
	// for _, historyMap := range tickerHistory {
	// 	for _, item := range historyMap {
	// 		fmt.Printf("%+v\t", item)
	// 	}
	// }

	//Loop through until we reach today
	for d := startDate; d.Before(today); d = d.AddDate(0, 0, 1) {
		var dayStr = d.Format("2006-01-02")
		var stockValue float64

		var lastTransactionOfDay *models.PortfolioEntity = nil
		for key := range portfolioEntities {
			//assetArr is the array of all transactions for that ticker
			assetArr := portfolioEntities[key]
			//find the amount of shares I had on dayStr
			//by making each the date the last instant of every day, we ensure that we count all of my trades made for the whole day
			lastTransaction := binarySearchEntityByDate(assetArr, d.Add(24*time.Hour-time.Nanosecond))

			//find price of stock on that day (specifially) then append it to stock value
			dayClose := tickerHistory[key][dayStr].Price
			stockValue += float64(lastTransaction.Shares) * float64(dayClose)

			if lastTransactionOfDay == nil {
				lastTransactionOfDay = lastTransaction
			} else if lastTransaction.CreatedAt.After(lastTransactionOfDay.CreatedAt) {
				lastTransactionOfDay = lastTransaction
			}
		}
		//append stockvalue to and date to array
		graphData = append(graphData, &dtos.PortfolioGraphDto{
			Date:        dayStr,
			StockAssets: stockValue,
			Cash:        lastTransactionOfDay.Cash,
		})
	}
	//todo: the graph goes up untill yesterday, final point of refernce in graph will just be the current portfolio value?

	return graphData, nil
}

func binarySearchEntityByDate(portfolio []*models.PortfolioEntity, date time.Time) *models.PortfolioEntity {
	left := 0
	right := len(portfolio) - 1
	for left < right {
		mid := (left + right + 1) / 2
		midDate := portfolio[mid].CreatedAt

		if midDate.After(date) {
			right = mid - 1
		} else {
			left = mid
		}
	}

	var targetEntity *models.PortfolioEntity = portfolio[left]
	if targetEntity.CreatedAt.After(date) {
		return &models.PortfolioEntity{
			Cash:   100000,
			Shares: 0,
		}
	}
	return targetEntity
}
