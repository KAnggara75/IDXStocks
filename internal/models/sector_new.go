package models

type SectorNew struct {
	Id   int
	Name string
}

type SubSector struct {
	Id       int
	Name     string
	SectorId int
}

type SectorSyncNewResponse struct {
	Sectors    []BasicResponse `json:"sectors"`
	SubSectors []BasicResponse `json:"sub_sectors"`
}
