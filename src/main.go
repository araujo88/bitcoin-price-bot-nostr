package main

import "log"

const URL = "https://rest.coinapi.io/v1/exchangerate/BTC/USD"

var API_KEY = goDotEnvVariable("API_KEY")

func main() {

	err := doPost(getRate())

	if err != nil {
		log.Fatal(err)
	}
}
