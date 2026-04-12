package models

type PasardanaSector struct {
	Id          int    `json:"Id"`
	Name        string `json:"Name"`
	NameEn      string `json:"NameEn"`
	Code        string `json:"Code"`
	Description string `json:"Description"`
}

type SectorResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
