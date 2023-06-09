package main

import (
	"fmt"
	"os"

	fiber "github.com/gofiber/fiber/v2"
)

type TickerResponse struct {
	Tickers []Ticker `json:"tickers"`
}

func main() {
	go runWebsocket()

	webApp := fiber.New()

	webApp.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	webApp.Get("/orderbook/:product", func(c *fiber.Ctx) error {
		product := c.Params("product")

		fmt.Println("product: ", product)
		fmt.Println("tickers: ", tickers)
		tickerList, ok := tickers[product]
		if !ok {
			return c.Status(404).SendString("No ticker found for product")
		}

		return c.JSON(TickerResponse{Tickers: tickerList})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	webApp.Listen("0.0.0.0:" + port)
}
