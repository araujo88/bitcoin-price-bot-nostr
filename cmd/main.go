package main

import (
	"log"
	"time"

	nostrbot "github.com/araujo88/bitcoin-price-bot-nostr/pkg"
)

func main() {

	for {
		err := nostrbot.DoPost()

		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Hour)
	}
}
