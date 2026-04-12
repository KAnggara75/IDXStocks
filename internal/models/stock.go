package models

type Stock struct {
	Id                 *int     `json:"id,omitempty"`
	Code               string   `json:"code"`
	CompanyName        string   `json:"company_name"`
	ListingDate        *string  `json:"listing_date,omitempty"`
	DelistingDate      *string  `json:"delisting_date,omitempty"`
	ListingBoard       string   `json:"listing_board"`
	Shares             int64    `json:"shares"`
	TotalEmployees     *string  `json:"total_employees,omitempty"`
	AnnualDividend     *float64 `json:"annual_dividend,omitempty"`
	GeneralInformation *string  `json:"general_information,omitempty"`
	FoundingDate       *string  `json:"founding_date,omitempty"`
	SectorId           *int     `json:"sector_id,omitempty"`
	SubSectorId        *int     `json:"sub_sector_id,omitempty"`
	IndustryId         *int     `json:"industry_id,omitempty"`
	SubIndustryId      *int     `json:"sub_industry_id,omitempty"`
}

type PasardanaStockDetail struct {
	Id                 int      `json:"Id"`
	Code               string   `json:"Code"`
	Name               string   `json:"Name"`
	TotalEmployees     *string  `json:"TotalEmployees"`
	ListingDate        *string  `json:"ListingDate"`
	AnnualDividend     *float64 `json:"AnnualDividend"`
	GeneralInformation *string  `json:"GeneralInformation"`
	FoundingDate       *string  `json:"FoundingDate"`
	FkNewSectorId      *int     `json:"fkNewSectorId"`
	FkNewSubSectorId   *int     `json:"fkNewSubSectorId"`
	FkNewIndustryId    *int     `json:"fkNewIndustryId"`
	FkNewSubIndustryId *int     `json:"fkNewSubIndustryId"`
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
