package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/coinpaprika/coinpaprika-api-go-client/v2/coinpaprika"
)

type csvTicker struct {
	tickerID     string
	lastUpdated  time.Time
	now          time.Time
	price        float64
	priceChanged bool
}

func (t csvTicker) record() []string {
	return []string{
		t.tickerID,
		t.lastUpdated.String(),
		fmt.Sprintf("%d", t.lastUpdated.Unix()),
		t.now.String(),
		fmt.Sprintf("%d", t.now.Unix()),
		fmt.Sprintf("%f", t.price),
		fmt.Sprintf("%t", t.priceChanged),
	}
}

func main() {
	apiKey := flag.String("key", "", "Coinpaprika API key")
	fileName := flag.String("file", "result.csv", "output file")
	tickerID := flag.String("ticker", "btc-bitcoin", "Coinpaprika ticker ID")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("API key is required")
	}

	client := coinpaprika.NewClient(nil, coinpaprika.WithAPIKey(*apiKey))

	file, err := os.Create(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	multiWriter := io.MultiWriter(os.Stdout, file)

	go func() {
		csvWriter := csv.NewWriter(multiWriter)
		_ = csvWriter.Write([]string{
			"ticker_id",                            // ticker ID from Coinpaprika API
			"last_update", "last_update_timestamp", // last_update field from Coinpaprika API
			"now", "now_timestamp", // current time for saving in CSV
			"price",         // price field from Coinpaprika API
			"price_changed", // if price changed since last update
		})

		var lastPrice float64
		timeTicker := time.NewTicker(10 * time.Second)
		for ; ; <-timeTicker.C {
			// gets https://api-pro.coinpaprika.com/v1/tickers/ every 10 seconds
			tickers, err := client.Tickers.List(nil)
			if err != nil {
				log.Fatal(err)
			}

			var ticker csvTicker
			for _, t := range tickers {
				if *t.ID == *tickerID {
					ticker.tickerID = *t.ID
					ticker.price = *t.Quotes["USD"].Price
					ticker.now = time.Now().UTC()
					lastUpdatedCoinpaprika, err := time.Parse(time.RFC3339, *t.LastUpdated)
					if err != nil {
						log.Fatal(err)
					}
					ticker.lastUpdated = lastUpdatedCoinpaprika

					if ticker.price != lastPrice {
						if lastPrice != 0 {
							ticker.priceChanged = true
						}
						lastPrice = ticker.price
					}
				}
			}

			_ = csvWriter.Write(ticker.record())
			csvWriter.Flush()
		}
	}()

	select {}
}
