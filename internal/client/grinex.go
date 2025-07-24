package client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"usdt/internal/models"
)

type GrinexClient struct {
	URL string
}

type depthItem struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

type depthResponse struct {
	Asks []depthItem `json:"asks"`
	Bids []depthItem `json:"bids"`
}

// GetDepth fetches order book depth from Grinex API and parses the best ask and bid prices.
func (c *GrinexClient) GetDepth() (*models.Rate, error) {
	resp, err := http.Get(c.URL + "?market=usdtrub")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}()
	var data depthResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data.Asks) == 0 || len(data.Bids) == 0 {
		return nil, errors.New("empty depth")
	}

	ask, err := strconv.ParseFloat(data.Asks[0].Price, 64)
	if err != nil {
		return nil, err
	}
	bid, err := strconv.ParseFloat(data.Bids[0].Price, 64)
	if err != nil {
		return nil, err
	}

	return &models.Rate{
		Ask:       ask,
		Bid:       bid,
		Timestamp: time.Now().UTC(),
	}, nil
}
