package services

import (
	"DellTradingApi/dtos"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func GetQuote(ticker string) (*dtos.StockResponseDto, error) {
	params := url.Values{}
	params.Add("symbol", ticker)

	resp, err := sendRequest(params, "quote")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rawJSON, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)

	}

	// var stockResponse *dtos.StockResponseDto = &dtos.StockResponseDto{}
	var stockResponse dtos.StockResponseDto
	unMarshalErr := json.Unmarshal(rawJSON, &stockResponse)
	if unMarshalErr != nil {
		return nil, unMarshalErr
	}

	return &stockResponse, nil
}

func sendRequest(params url.Values, endpoint string) (*http.Response, error) {
	urlStr := fmt.Sprintf("https://api.twelvedata.com/%s", endpoint)
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	client := &http.Client{}
	params.Add("apikey", os.Getenv("STOCK_API_KEY"))
	req.URL.RawQuery = params.Encode()

	// Send the request and get a response
	resp, err := client.Do(req)
	return resp, err
}
