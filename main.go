package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	ChannelTicker    string = "ticker"
	Orderbookchannel string = "level2"
	SubscriptionType string = "subscribe"
	Apikey           string = "test"
	secret           string = "secret"
)

const Address = "wss://advanced-trade-ws.coinbase.com"

func buildSing(timeNow int64, products []string) string {
	stringToSign := fmt.Sprintln(timeNow, Orderbookchannel, strings.Join(products, ","))
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return hex.EncodeToString(h.Sum(nil))
}

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
	timeNow := time.Now().Unix()
	productIds := []string{"BTC-USD", "ETH-USD"}
	signature := buildSing(timeNow, productIds)

	sub := Message{
		Type:       SubscriptionType,
		ProductIDs: productIds,
		Channel:    Orderbookchannel,
		Timestamp:  timeNow,
		ApiKey:     Apikey,
		Signature:  signature,
	}

	if err := wsjson.Write(context.Background(), conn, sub); err != nil {
		log.Fatal(err)
		return
	}

	// Listen indefinitely for messages until an interrupt signal is received.
	for {
		var orderbook OrderBook
		if err := wsjson.Read(context.Background(), conn, &orderbook); err != nil {
			log.Fatal(err)
			return
		}

		//fmt.Printf("%s: %s | %s | %s | %s \n", ticker.ProductID, ticker.Price, ticker.BestBid, ticker.BestAsk, ticker.Time)
		fmt.Println(orderbook)
	}
}
