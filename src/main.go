package main

import (
	"log"
	"time"
)

func main() {

	for {
		err := doPost()

		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(1 * time.Hour)
	}
}
