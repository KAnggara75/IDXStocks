package models

type PasardanaSearchResult struct {
	Id                 int    `json:"Id"`
	Name               string `json:"Name"`
	Code               string `json:"Code"`
	NewSubIndustryId   int    `json:"NewSubIndustryId"`
	NewSubIndustryName string `json:"NewSubIndustryName"`
	NewIndustryId      int    `json:"NewIndustryId"`
	NewIndustryName    string `json:"NewIndustryName"`
	NewSubSectorId     int    `json:"NewSubSectorId"`
	NewSubSectorName   string `json:"NewSubSectorName"`
	NewSectorId        int    `json:"NewSectorId"`
	NewSectorName      string `json:"NewSectorName"`
}

type Industry struct {
	Id   int
	Name string
}

type SubIndustry struct {
	Id         int
	Name       string
	IndustryId int
}

type IndustrySyncResponse struct {
	Industries    []BasicResponse `json:"industries"`
	SubIndustries []BasicResponse `json:"sub_industries"`
}

type BasicResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
