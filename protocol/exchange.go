package protocol

type BithumbTickerData struct {
	Last  float64 `json:"closing_price,string"`
	First float64 `json:"opening_price,string"`
	// MinPrice     float64 `json:"min_price,string"`
	// MaxPrice     float64 `json:"max_price,string"`
	// AveragePrice float64 `json:"average_price,string"`
	// UnitsTraded  float64 `json:"units_traded,string"`
	// Volume1day   float64 `json:"volume_1day,string"`
	// Volume7day   float64 `json:"volume_7day,string"`
	// BuyPrice     float64 `json:"buy_price,string"`
	// SellPrice    float64 `json:"sell_price,string"`
	// Date string `json:"date"`
}

type BithumbTickerResponse struct {
	Status string            `json:"status"`
	Data   BithumbTickerData `json:"data"`
}

type BithumbTickerAllResponse struct {
	Status string                 `json:"status"`
	Data   map[string]interface{} `json:"data"`
}

type CoinoneTickerResponse struct {
	Last  float64 `json:"last,string,string"`
	First float64 `json:"first,string"`
	// Low       float64 `json:"lowmin_price,string"`
	// High      float64 `json:"high,string"`
	// Volume    float64 `json:"volume,string"`
	// Currency  float64 `json:"currency,string"`
	Timestamp string `json:"timestamp"`
	Result    string `json:"result"`
	ErrorCode string `json:"errorCode"`
}

type KorbitTickerData struct {
	Timestamp int     `json:"timestamp"`
	Last      float64 `json:"last,string"`
	// Bid       float64 `json:"bid,string"`
	// Ask       float64 `json:"ask,string"`
	// Low       float64 `json:"low,string"`
	// High      float64 `json:"high,string"`
	// Volume    float64 `json:"volume,string"`
	Change        float64 `json:"change,string"`
	ChangePercent float64 `json:"changePercent,string"`
}
