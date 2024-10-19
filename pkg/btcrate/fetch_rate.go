package btcrate

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/config"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/message"
)

const BASE_URL = "https://rest.coinapi.io/v1/exchangerate/BTC/"

var API_KEY = config.GetDotEnvVariable("API_KEY")

func FetchRate(currency string) float64 {

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, BASE_URL+currency, nil)
	var response []byte

	req.Header.Set("X-CoinAPI-Key", API_KEY)

	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		response, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		fmt.Printf("Error: status code %d", res.StatusCode)
	}

	message := message.Message{}

	err = json.Unmarshal(response, &message)

	if err != nil {
		log.Fatal(err)
	}

	return message.Rate
}
