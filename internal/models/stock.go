package models

type Stock struct {
	Code          string `json:"Code"`
	CompanyName   string `json:"CompanyName"`
	ListingDate   string `json:"ListingDate"`
	DelistingDate string `json:"DelistingDate"`
	ListingBoard  string `json:"ListingBoard"`
	Shares        int64  `json:"Shares"`
}
