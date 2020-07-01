package api

type ErrorResponse struct {
	Error string `json:"Error"`
}

type RankRequest struct {
	Limit int    `url:"limit"`
	Page  int    `url:"page"`
	TSYM  string `url:"tsym"`
}

type RankResponse struct {
	Rank     int    `json:"Rank"`
	Name     string `json:"Name"`
	FullName string `json:"FullName"`
}

type PriceRequest struct {
	Symbol  []string `url:"symbol"`
	Convert string   `url:"convert"`
}

type PriceResponse map[string]map[string]float32

type Data struct {
	Rank     int     `json:"Rank"`
	Name     string  `json:"Name"`
	FullName string  `json:"FullName"`
	Price    float32 `json:"Price"`
	Currency string  `json:"Currency"`
}
