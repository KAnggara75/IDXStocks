package models

type Stock struct {
	Id            *int   `json:"id,omitempty"`
	Code          string `json:"code"`
	CompanyName   string `json:"company_name"`
	ListingDate   string `json:"listing_date"`
	DelistingDate string `json:"delisting_date,omitempty"`
	ListingBoard  string `json:"listing_board"`
	Shares        int64  `json:"shares"`
}

type PasardanaStock struct {
	Id   int    `json:"Id"`
	Code string `json:"Code"`
	Name string `json:"Name"`
}

type StockResponse struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
