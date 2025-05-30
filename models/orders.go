package models

import "time"

type Order struct {
	ID           int64     `json:"id"`
	Symbol       string    `json:"symbol"`
	Side         OrderSide `json:"side"`
	Type         OrderType `json:"type"`
	Price        float64   `json:"price,omitempty"`
	Quantity     int       `json:"quantity"`
	RemainingQty int       `json:"remaining_quantity"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type Trade struct {
	ID          int64     `json:"id"`
	BuyOrderID  int64     `json:"buy_order_id"`
	SellOrderID int64     `json:"sell_order_id"`
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrderType string
type OrderSide string

const (
	LimitOrder  OrderType = "limit"
	MarketOrder OrderType = "market"

	BuySide  OrderSide = "buy"
	SellSide OrderSide = "sell"
)
