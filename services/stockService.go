package services

import (
	"DellTradingApi/dtos"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func GetQuote(ticker string) (*dtos.StockQuoteDto, error) {
	params := url.Values{}
	params.Add("symbol", ticker)

	resp, err := sendRequest(params, "quote")
	//todo: put all of this in sendRequest function?
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var stockResponse dtos.StockQuoteDto

	unMarshalErr := json.Unmarshal(rawJSON, &stockResponse)
	if unMarshalErr != nil {
		return nil, unMarshalErr
	}

	if stockResponse.Name == "" {
		return nil, fmt.Errorf("invalid ticker")
	}
	return &stockResponse, nil
}

func GetHistory(ticker string, startDate string) (map[string]*dtos.TimeSeriesQuoteDto, error) {
	params := url.Values{}
	params.Add("symbol", ticker)
	params.Add("interval", "1day")
	params.Add("start_date", startDate)

	resp, err := sendRequest(params, "time_series")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var stockHistory dtos.TimeSeriesDto
	unMarshalErr := json.Unmarshal(rawJSON, &stockHistory)
	if unMarshalErr != nil {
		return nil, unMarshalErr
	}

	priceHistory := stockHistory.Values
	historyMap := make(map[string]*dtos.TimeSeriesQuoteDto)

	nextDate := priceHistory[len(priceHistory)-1].Date
	lastQuote := priceHistory[len(priceHistory)-1]
	for i := len(priceHistory) - 1; i > -1; {
		timeSlice := priceHistory[i]
		dateKey := strings.Split(timeSlice.Date, " ")[0]
		if timeSlice.Date != nextDate {
			dateKey = strings.Split(nextDate, " ")[0]
			historyMap[dateKey] = lastQuote
		} else {
			historyMap[dateKey] = timeSlice
			lastQuote = timeSlice
			i--
		}

		parsed, err := time.Parse("2006-01-02", nextDate)
		if err != nil {
			return nil, err
		}
		nextDate = parsed.AddDate(0, 0, 1).Format("2006-01-02")

	}

	return historyMap, nil
}

func sendRequest(params url.Values, endpoint string) (*http.Response, error) {
	urlStr := fmt.Sprintf("https://api.twelvedata.com/%s", endpoint)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	client := &http.Client{}
	params.Add("apikey", os.Getenv("STOCK_API_SECRET"))

	req.URL.RawQuery = params.Encode()

	// Send the request and get a response
	resp, err := client.Do(req)
	//try using .GET INSTEAD
	return resp, err
}
