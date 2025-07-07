package main

import (
	"fmt"
	"math/rand"
	"time"
)

// RandomDataGenerator generates random stock data for simulation
type RandomDataGenerator struct {
	stockCodes  []string
	stockPrices []int
	buyOrSell   []string
	rand        *rand.Rand
}

// NewRandomDataGenerator creates a new instance of RandomDataGenerator with seed
func NewRandomDataGenerator() *RandomDataGenerator {
	source := rand.NewSource(time.Now().UnixNano())
	return &RandomDataGenerator{
		stockCodes:  []string{"AAPL", "GOOGL", "MSFT", "AMZN", "TSLA", "META", "NFLX", "INTC", "AMD", "NVDA"},
		stockPrices: []int{100, 200, 300, 400, 500, 150, 250, 350, 450, 550},
		buyOrSell:   []string{"B", "S"},
		rand:        rand.New(source),
	}
}

// GetData returns a random stock data entry in the format "StockCode|Price|B/S|Amount"
func (g *RandomDataGenerator) GetData() string {
	stockCode := g.stockCodes[g.rand.Intn(len(g.stockCodes))]
	stockPrice := g.stockPrices[g.rand.Intn(len(g.stockPrices))]
	tradeType := g.buyOrSell[g.rand.Intn(len(g.buyOrSell))]
	quantity := g.rand.Intn(100) + 1
	
	return fmt.Sprintf("%s|%d|%s|%d", stockCode, stockPrice, tradeType, quantity)
}
