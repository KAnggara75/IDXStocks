package models

type PasardanaSearchResult struct {
	NewSubIndustryId   int    `json:"NewSubIndustryId"`
	NewSubIndustryName string `json:"NewSubIndustryName"`
	NewIndustryId      int    `json:"NewIndustryId"`
	NewIndustryName    string `json:"NewIndustryName"`
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
