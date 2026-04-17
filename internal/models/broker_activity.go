package models

import "time"

type BrokerActivity struct {
	BrokerCode string    `json:"broker_code" db:"broker_code"`
	StockCode  string    `json:"stock_code" db:"stock_code"`
	Date       time.Time `json:"date" db:"date"`
	Side       string    `json:"side" db:"side"` // "buy" or "sell"
	Lot        int64     `json:"lot" db:"lot"`
	Value      int64     `json:"value" db:"value"`
	AvgPrice   float64   `json:"avg_price" db:"avg_price"`
	Freq       int64     `json:"freq" db:"freq"`
}

type ExodusBrokerActivityItem struct {
	StockCode  string  `json:"stock_code"`
	BrokerCode string  `json:"broker_code"`
	Date       string  `json:"date"`
	Value      float64 `json:"value"`
	Lot        float64 `json:"lot"`
	AvgPrice   float64 `json:"avg_price"`
	Freq       int64   `json:"freq"`
}

type ExodusBrokerActivityResponse struct {
	Message string `json:"message"`
	Data    struct {
		BrokerActivityTransaction struct {
			BrokersBuy  []ExodusBrokerActivityItem `json:"brokers_buy"`
			BrokersSell []ExodusBrokerActivityItem `json:"brokers_sell"`
		} `json:"broker_activity_transaction"`
	} `json:"data"`
}

type SyncBrokerActivityParams struct {
	BrokerCode      string `query:"broker_code"`
	From            string `query:"from"`
	To              string `query:"to"`
	TransactionType string `query:"transaction_type"`
	MarketBoard     string `query:"market_board"`
	InvestorType    string `query:"investor_type"`
}
