package main

import (
	"log"
	"time"

	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/post"
)

func main() {

	for {
		err := post.Post()

		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(24 * time.Hour)
	}
}
