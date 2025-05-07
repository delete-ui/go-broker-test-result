package model

import "time"

type Trade struct {
	Id        int64     `json:"id"`
	Account   string    `json:"account"`
	Symbol    string    `json:"symbol"`
	Volume    float64   `json:"volume"`
	Open      float64   `json:"open"`
	Close     float64   `json:"close"`
	Side      string    `json:"side"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountStats struct {
	Account string  `json:"account"`
	Trades  float64 `json:"trades"`
	Profit  float64 `json:"profit"`
}

type TradeRequest struct {
	Account string  `json:"account"`
	Symbol  string  `json:"symbol"`
	Volume  float64 `json:"volume"`
	Open    float64 `json:"open"`
	Close   float64 `json:"close"`
	Side    string  `json:"side"`
}
