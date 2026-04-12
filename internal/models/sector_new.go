package models

type SectorNew struct {
	Id          int
	Code        *string
	Name        string
	NameEn      *string
	Description *string
}

type SubSector struct {
	Id          int
	SectorId    int
	Code        *string
	Name        string
	NameEn      *string
	Description *string
}

type PasardanaNewSector struct {
	Id          int     `json:"Id"`
	Code        *string `json:"Code"`
	Name        string  `json:"Name"`
	NameEn      *string `json:"NameEn"`
	Description *string `json:"Description"`
}

type PasardanaNewSubSector struct {
	Id            int     `json:"Id"`
	FkNewSectorId int     `json:"fkNewSectorId"`
	Code          *string `json:"Code"`
	Name          string  `json:"Name"`
	NameEn        *string `json:"NameEn"`
	Description   *string `json:"Description"`
}

type BasicResponseWithCode struct {
	Id   int     `json:"id"`
	Code *string `json:"code,omitempty"`
	Name string  `json:"name"`
}

type SectorSyncNewResponse struct {
	Sectors    []BasicResponseWithCode `json:"sectors"`
	SubSectors []BasicResponseWithCode `json:"sub_sectors"`
}
