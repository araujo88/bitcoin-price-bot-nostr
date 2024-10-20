package coinapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/config"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/responses"
)

const BASE_URL = "https://rest.coinapi.io/v1/"

var API_KEY = config.GetDotEnvVariable("API_KEY")

// makeRequest handles the common operations of making an HTTP request to the CoinAPI
func makeRequest(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("X-CoinAPI-Key", API_KEY)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body) // Read the body
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, error reading body: %v", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return response, nil
}

// FetchRate retrieves the current rate of Bitcoin in the specified currency
func FetchRate(currency string) (float64, error) {
	url := BASE_URL + "exchangerate/BTC/" + currency
	response, err := makeRequest(url)
	if err != nil {
		return 0, err
	}

	var message responses.CurrencyRate
	if err := json.Unmarshal(response, &message); err != nil {
		return 0, fmt.Errorf("unmarshal response: %w", err)
	}

	return message.Rate, nil
}

// FetchDailyVariation retrieves the daily variation in percentage of the Bitcoin price for a specified currency
func FetchDailyVariation(currency string) (float64, error) {
	symbolID := fmt.Sprintf("BITSTAMP_SPOT_BTC_%s", currency)               // Adjust the exchange as needed
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02T00:00:00") // Format yesterday's date in ISO 8601 format
	url := fmt.Sprintf("%sohlcv/%s/history?period_id=1DAY&time_start=%s&limit=1", BASE_URL, symbolID, yesterday)

	response, err := makeRequest(url)
	if err != nil {
		return 0, err
	}

	var data []responses.OHLCVData
	if err := json.Unmarshal(response, &data); err != nil {
		return 0, fmt.Errorf("error unmarshalling OHLCV data: %w", err)
	}

	if len(data) > 0 {
		ohlcv := data[0]
		variation := ((ohlcv.PriceClose - ohlcv.PriceOpen) / ohlcv.PriceOpen) * 100
		return variation, nil
	}

	return 0, fmt.Errorf("no data found")
}
