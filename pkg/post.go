package nostrbot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func DoPost() error {
	profile := ""
	cfg, err := loadConfig(profile)

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

	rate_usd := 1 / getRate("USD") / 0.00000001
	rate_eur := 1 / getRate("EUR") / 0.00000001
	rate_jpy := 1 / getRate("JPY") / 0.00000001
	rate_gbp := 1 / getRate("GBP") / 0.00000001
	rate_brl := 1 / getRate("BRL") / 0.00000001

	price_string := fmt.Sprintf("1 USD = %.0f sats\n1 EUR = %0.f sats\n1 JPY = %0.f sats\n1 GBP = %0.f sats\n1 BRL = %0.f sats", rate_usd, rate_eur, rate_jpy, rate_gbp, rate_brl)
	ev.Content = price_string

	ev.CreatedAt = time.Now()
	ev.Kind = nostr.KindTextNote
	ev.Sign(sk)

	var success atomic.Int64
	cfg.Do(Relay{Write: true}, func(relay *nostr.Relay) {
		status := relay.Publish(context.Background(), ev)
		if cfg.verbose {
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
