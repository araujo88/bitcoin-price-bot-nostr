# Bitcoin Price Bot for Nostr

[![license](https://img.shields.io/badge/license-GPL-green)](https://raw.githubusercontent.com/araujo88/bitcoin-price-bot-nostr/main/LICENSE)
[![build](https://github.com/araujo88/bitcoin-price-bot-nostr/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/araujo88/bitcoin-price-bot-nostr/actions/workflows/go.yml)

This bot posts Bitcoin price information to Nostr relays at regular intervals.

## Setup

1. Clone the repository:

   ```
   git clone https://github.com/araujo88/bitcoin-price-bot-nostr.git
   cd bitcoin-price-bot-nostr
   ```

2. Create a configuration file:
   Create a file named `config.json` in the directory `/home/ubuntu/.config/bitcoin-price-bot/` with the following structure:

   ```json
   {
     "relays": {
       "wss://relay1.example.com": {
         "read": true,
         "write": true,
         "search": false
       },
       "wss://relay2.example.com": {
         "read": true,
         "write": true,
         "search": false
       }
     },
     "privatekey": "your_nsec_private_key_here"
   }
   ```

   Replace the relay URLs with the Nostr relays you want to use, and add your nsec (private key) to the "privatekey" field.

3. Set up environment variables:
   Create a `.env` file in the root directory of the project and add your CoinAPI key:

   ```
   API_KEY=your_coinapi_key_here
   ```

## Usage

To run the bot:

```
go run main.go
```

The bot will fetch Bitcoin prices for USD, EUR and BRL daily and post them to the configured Nostr relays.

## Configuration

- The `config.json` file specifies the Nostr relays to connect to and the private key for signing messages.
- The `.env` file contains the API key for CoinAPI, which is used to fetch Bitcoin price data.

## Features

- Fetches Bitcoin prices for multiple currencies (USD, EUR, JPY, GBP, BRL)
- Posts price information to Nostr relays hourly
- Configurable relays and private key

## Dependencies

This project uses the following main dependencies:

- `github.com/nbd-wtf/go-nostr`: For Nostr protocol implementation
- `github.com/joho/godotenv`: For loading environment variables

Make sure to run `go mod tidy` to fetch all required dependencies.

## Note

Ensure that you keep your private key (`nsec`) secure and do not share it publicly. The `config.json` file containing your private key should be kept confidential.
