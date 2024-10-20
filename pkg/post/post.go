package post

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/coinapi"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/config"
	"github.com/araujo88/bitcoin-price-bot-nostr/pkg/relay"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

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
		return errors.New("error fetching rate for USD")
	}
	rate_eur, err := coinapi.FetchRate("EUR")
	if err != nil {
		return errors.New("error fetching rate for EUR")
	}
	rate_brl, err := coinapi.FetchRate("BRL")
	if err != nil {
		return errors.New("error fetching rate for BRL")
	}

	daily_variation_usd, err := coinapi.FetchDailyVariation("USD")
	if err != nil {
		return errors.New("error fetching daily vartion for USD")
	}

	daily_variation_eur, err := coinapi.FetchDailyVariation("EUR")
	if err != nil {
		return errors.New("error fetching daily vartion for EUR")
	}

	daily_variation_brl, err := coinapi.FetchDailyVariation("BRL")
	if err != nil {
		return errors.New("error fetching daily vartion for BRL")
	}

	content := fmt.Sprintf(`1 BTC = %.0f USD (%.2f %%)\n
	1 BTC = %0.f EUR (%.2f %%)\n
	1 BTC = %0.f BRL (%.2f %%)\n`,
		rate_usd, daily_variation_usd,
		rate_eur, daily_variation_eur,
		rate_brl, daily_variation_brl)

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
