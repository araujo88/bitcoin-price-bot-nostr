package post

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/coinapi"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/config"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/relay"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func formatWithSeparator(value float64) string {
	// Convert number to string without decimals
	str := strconv.FormatFloat(value, 'f', 0, 64)
	if len(str) <= 3 {
		return str
	}

	// Reverse the string to start placing commas every three digits
	n := len(str)
	var buffer bytes.Buffer
	for i := n - 1; i >= 0; i-- {
		buffer.WriteByte(str[i])
		if (n-i)%3 == 0 && i != 0 {
			buffer.WriteByte(',')
		}
	}

	// Reverse the string back to the original order
	runes := bytes.Runes(buffer.Bytes())
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func Post() error {
	profile := ""
	cfg, err := config.LoadConfig(profile)

	if err != nil {
		return err
	}

	var sk string

	if _, s, err := nip19.Decode(cfg.PrivateKey); err != nil {
		return err
	} else {
		sk = s.(string)
	}
	ev := nostr.Event{}
	if pub, err := nostr.GetPublicKey(sk); err == nil {
		if _, err := nip19.EncodePublicKey(pub); err != nil {
			return err
		}
		ev.PubKey = pub
	} else {
		return err
	}

	rate_usd, err := coinapi.FetchRate("USD")
	if err != nil {
		return err
	}
	rate_eur, err := coinapi.FetchRate("EUR")
	if err != nil {
		return err
	}
	rate_brl, err := coinapi.FetchRate("BRL")
	if err != nil {
		return err
	}

	daily_variation_usd, err := coinapi.FetchDailyVariation("BITSTAMP", "USD")
	if err != nil {
		return err
	}

	daily_variation_eur, err := coinapi.FetchDailyVariation("BITSTAMP", "EUR")
	if err != nil {
		return err
	}

	daily_variation_brl, err := coinapi.FetchDailyVariation("BINANCE", "BRL")
	if err != nil {
		return err
	}

	content := fmt.Sprintf(`1 BTC = %s USD (%.2f %%)
	1 BTC = %s EUR (%.2f %%)
	1 BTC = %s BRL (%.2f %%)`,
		formatWithSeparator(rate_usd), daily_variation_usd,
		formatWithSeparator(rate_eur), daily_variation_eur,
		formatWithSeparator(rate_brl), daily_variation_brl)

	ev.Content = content

	ev.CreatedAt = time.Now()
	ev.Kind = nostr.KindTextNote
	ev.Sign(sk)

	var success atomic.Int64
	cfg.Do(relay.Relay{Write: true}, func(relay *nostr.Relay) {
		status := relay.Publish(context.Background(), ev)
		if cfg.Verbose {
			fmt.Fprintln(os.Stderr, relay.URL, status)
		}
		if status != nostr.PublishStatusFailed {
			success.Add(1)
		}
	})
	if success.Load() == 0 {
		return errors.New("cannot post")
	}
	return nil
}
