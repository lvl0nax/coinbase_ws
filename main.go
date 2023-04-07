package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	ChannelTicker    string = "ticker"
	SubscriptionType string = "subscribe"
)

type Message struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Ticker struct {
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	ProductID string    `json:"product_id"`
	Price     string    `json:"price"`
	Open24H   string    `json:"open_24h"`
	BestBid   string    `json:"best_bid"`
	BestAsk   string    `json:"best_ask"`
}

const Address = "wss://ws-feed.exchange.coinbase.com"

func main() {
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt)

	u, err := url.Parse(Address)
	if err != nil {
		log.Fatal(err)
	}

	conn, _, err := websocket.Dial(context.Background(), u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(websocket.StatusInternalError, "the sky is falling")

	// Subscribe to the BTC-USD ticker channel.
	sub := Message{
		Type:       SubscriptionType,
		ProductIDs: []string{"BTC-USD"},
		Channels:   []string{ChannelTicker},
	}

	if err := wsjson.Write(context.Background(), conn, sub); err != nil {
		log.Fatal(err)
		return
	}

	// Listen indefinitely for messages until an interrupt signal is received.
	for {
		var ticker Ticker
		if err := wsjson.Read(context.Background(), conn, &ticker); err != nil {
			log.Fatal(err)
			return
		}

		fmt.Printf("%s: %s | %s | %s | %s \n", ticker.ProductID, ticker.Price, ticker.BestBid, ticker.BestAsk, ticker.Time)
	}
}
