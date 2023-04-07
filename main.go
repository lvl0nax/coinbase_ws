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

// Message is the message sent to the Coinbase websocket advanced API.
//
//	{
//		"type": "subscribe",
//		"product_ids": [
//			"ETH-USD",
//			"BTC-USD"
//		],
//		"channel": "ticker",
//		"signature": "XYZ",
//		"api_key": "XXX",
//		"timestamp": 1675974199
//	}
type Message struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channel    string   `json:"channel"`
	Timestamp  int64    `json:"timestamp"`
	ApiKey     string   `json:"api_key"`
	Signature  string   `json:"signature"`
}

// Ticker messsage
//
//	{
//		"channel": "ticker",
//		"client_id": "",
//		"timestamp": "2023-02-09T20:30:37.167359596Z",
//		"sequence_num": 0,
//		"events": [
//			{
//				"type": "snapshot",
//				"tickers": [
//					{
//						"type": "ticker",
//						"product_id": "BTC-USD",
//						"price": "21932.98",
//						"volume_24_h": "16038.28770938",
//						"low_24_h": "21835.29",
//						"high_24_h": "23011.18",
//						"low_52_w": "15460",
//						"high_52_w": "48240",
//						"price_percent_chg_24_h": "-4.15775596190603"
//					}
//				]
//			}
//		]
//	}
type Ticker struct {
	Channel     string    `json:"channel"`
	ClientID    string    `json:"client_id"`
	Timestamp   time.Time `json:"timestamp"`
	SequenceNum int       `json:"sequence_num"`
	Events      []struct {
		Type    string `json:"type"`
		Tickers []struct {
			Type      string `json:"type"`
			ProductID string `json:"product_id"`
			Price     string `json:"price"`
			Volume24H string `json:"volume_24_h"`
		}
	}
}

// OrderBook messsage
//
//	{
//		"channel": "l2_data",
//		"client_id": "",
//		"timestamp": "2023-02-09T20:32:50.714964855Z",
//		"sequence_num": 0,
//		"events": [
//			{
//				"type": "snapshot",
//				"product_id": "BTC-USD",
//				"updates": [
//					{
//						"side": "bid",
//						"event_time": "1970-01-01T00:00:00Z",
//						"price_level": "21921.73",
//						"new_quantity": "0.06317902"
//					},
//					{
//						"side": "bid",
//						"event_time": "1970-01-01T00:00:00Z",
//						"price_level": "21921.3",
//						"new_quantity": "0.02"
//					},
//				]
//			}
//		]
//	}
type OrderBook struct {
	Channel     string    `json:"channel"`
	ClientID    string    `json:"client_id"`
	Timestamp   time.Time `json:"timestamp"`
	SequenceNum int       `json:"sequence_num"`
	Events      []struct {
		Type      string `json:"type"`
		ProductID string `json:"product_id"`
		Updates   []struct {
			Side       string    `json:"side"`
			EventTime  time.Time `json:"event_time"`
			PriceLevel string    `json:"price_level"`
			NewQty     string    `json:"new_quantity"`
		}
	}
}

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
