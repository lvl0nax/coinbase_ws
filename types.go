package main

import (
	"time"
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
