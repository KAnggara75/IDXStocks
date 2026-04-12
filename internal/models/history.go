package models

import "time"

type StockHistory struct {
	Code                string    `json:"code"`
	Date                time.Time `json:"date"`
	Previous            *float64  `json:"previous"`
	OpenPrice           *float64  `json:"open_price"`
	FirstTrade          *float64  `json:"first_trade"`
	High                *float64  `json:"high"`
	Low                 *float64  `json:"low"`
	Close               *float64  `json:"close"`
	Change              *float64  `json:"change"`
	Volume              *float64  `json:"volume"`
	Value               *float64  `json:"value"`
	Frequency           *float64  `json:"frequency"`
	IndexIndividual     *float64  `json:"index_individual"`
	Offer               *float64  `json:"offer"`
	OfferVolume         *float64  `json:"offer_volume"`
	Bid                 *float64  `json:"bid"`
	BidVolume           *float64  `json:"bid_volume"`
	ListedShares        *float64  `json:"listed_shares"`
	TradebleShares      *float64  `json:"tradeble_shares"`
	WeightForIndex      *float64  `json:"weight_for_index"`
	ForeignSell         *float64  `json:"foreign_sell"`
	ForeignBuy          *float64  `json:"foreign_buy"`
	DelistingDate       *string   `json:"delisting_date"`
	NonRegularVolume    *float64  `json:"non_regular_volume"`
	NonRegularValue     *float64  `json:"non_regular_value"`
	NonRegularFrequency *float64  `json:"non_regular_frequency"`
	LastModified        time.Time `json:"last_modified"`
}

type SyncHistoryRequest struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type PasardanaHistoryResponse struct {
	Code                 string   `json:"Code"`
	PrevClosingPrice     *float64 `json:"PrevClosingPrice"`
	AdjustedClosingPrice *float64 `json:"AdjustedClosingPrice"`
	AdjustedOpenPrice    *float64 `json:"AdjustedOpenPrice"`
	AdjustedHighPrice    *float64 `json:"AdjustedHighPrice"`
	AdjustedLowPrice     *float64 `json:"AdjustedLowPrice"`
	Volume               *float64 `json:"Volume"`
	Frequency            *float64 `json:"Frequency"`
	Value                *float64 `json:"Value"`
	LastDate             *string  `json:"LastDate"`
}

type IdxSummaryData struct {
	Date                string   `json:"Date"`
	StockCode           string   `json:"StockCode"`
	Previous            *float64 `json:"Previous"`
	OpenPrice           *float64 `json:"OpenPrice"`
	FirstTrade          *float64 `json:"FirstTrade"`
	High                *float64 `json:"High"`
	Low                 *float64 `json:"Low"`
	Close               *float64 `json:"Close"`
	Change              *float64 `json:"Change"`
	Volume              *float64 `json:"Volume"`
	Value               *float64 `json:"Value"`
	Frequency           *float64 `json:"Frequency"`
	IndexIndividual     *float64 `json:"IndexIndividual"`
	Offer               *float64 `json:"Offer"`
	DelistingDate       string   `json:"DelistingDate"`
	OfferVolume         *float64 `json:"OfferVolume"`
	Bid                 *float64 `json:"Bid"`
	BidVolume           *float64 `json:"BidVolume"`
	ListedShares        *float64 `json:"ListedShares"`
	TradebleShares      *float64 `json:"TradebleShares"`
	WeightForIndex      *float64 `json:"WeightForIndex"`
	ForeignSell         *float64 `json:"ForeignSell"`
	ForeignBuy          *float64 `json:"ForeignBuy"`
	NonRegularVolume    *float64 `json:"NonRegularVolume"`
	NonRegularValue     *float64 `json:"NonRegularValue"`
	NonRegularFrequency *float64 `json:"NonRegularFrequency"`
}

type IdxSummaryResponse struct {
	Data []IdxSummaryData `json:"data"`
}
