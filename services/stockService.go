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

func GetHistory(ticker string, startDate time.Time) (map[string]*dtos.TimeSeriesQuoteDto, error) {
	params := url.Values{}
	params.Add("symbol", ticker)
	params.Add("interval", "1day")

	dayOfWeek := startDate.Weekday()
	startDateStr := startDate.Format("2006-01-02")
	var backOffset bool
	if dayOfWeek == time.Saturday {
		startDateStr = startDate.AddDate(0, 0, -1).Format("2006-01-02")
		backOffset = true
	}
	if dayOfWeek == time.Sunday {
		startDateStr = startDate.AddDate(0, 0, -2).Format("2006-01-02")
		backOffset = true
	}
	params.Add("start_date", startDateStr)

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

	now := time.Now()
	today := time.Date(
		now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
	)
	lastQuote := priceHistory[len(priceHistory)-1]
	dateKey := startDate
	for i := len(priceHistory) - 1; dateKey.Before(today) || i > -1; {
		dateStr := strings.Split(dateKey.Format("2006-01-02"), " ")[0]

		if i > -1 && dateStr == priceHistory[i].Date {
			timeSlice := priceHistory[i]
			//if they match, all is normal
			historyMap[dateStr] = timeSlice
			lastQuote = timeSlice
			i--
		} else {
			//set the value equal to the last value we had
			historyMap[dateStr] = lastQuote

			//if our stock api had to start before our start date to compensate for weekends
			//and we just insert our sunday, value, decrement the api response pointer so we can
			//reach the next value, and have our two date pointers sync up with one another
			if backOffset && dateKey.Weekday() == time.Sunday {
				i--

				//if the stock api for some reason contains data on a sunday
				//decrement it again, and go to monday
				if dateStr == priceHistory[i].Date {
					i--
				}
				backOffset = false
			}
		}
		dateKey = dateKey.AddDate(0, 0, 1)
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
