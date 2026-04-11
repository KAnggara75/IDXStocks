package models

type Stock struct {
	Code          string `json:"code"`
	CompanyName   string `json:"company_name"`
	ListingDate   string `json:"listing_date"`
	DelistingDate string `json:"delisting_date"`
	ListingBoard  string `json:"listing_board"`
	Shares        int64  `json:"shares"`
}
