package models

type Broker struct {
	Code           string `json:"code" db:"code"`
	Name           string `json:"name" db:"name"`
	InvestorType   string `json:"investor_type" db:"investor_type"`
	TotalValue     int64  `json:"total_value,string" db:"total_value"`
	NetValue       int64  `json:"net_value,string" db:"net_value"`
	BuyValue       int64  `json:"buy_value,string" db:"buy_value"`
	SellValue      int64  `json:"sell_value,string" db:"sell_value"`
	TotalVolume    int64  `json:"total_volume,string" db:"total_volume"`
	TotalFrequency int64  `json:"total_frequency,string" db:"total_frequency"`
	Group          string `json:"group" db:"broker_group"`
}

type BrokerSeederResponse struct {
	Message string `json:"message"`
	Data    struct {
		List []Broker `json:"list"`
	} `json:"data"`
}
