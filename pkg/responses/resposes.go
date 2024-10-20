package responses

type CurrencyRate struct {
	Time         string  `json:"time"`
	AssetIdBase  string  `json:"asset_id_base"`
	AssetIdQuote string  `json:"asset_id_quote"`
	Rate         float64 `json:"rate"`
}

type OHLCVData struct {
	TimePeriodStart string  `json:"time_period_start"`
	TimePeriodEnd   string  `json:"time_period_end"`
	TimeOpen        string  `json:"time_open"`
	TimeClose       string  `json:"time_close"`
	PriceOpen       float64 `json:"price_open"`
	PriceHigh       float64 `json:"price_high"`
	PriceLow        float64 `json:"price_low"`
	PriceClose      float64 `json:"price_close"`
	VolumeTraded    float64 `json:"volume_traded"`
	TradesCount     int     `json:"trades_count"`
}
