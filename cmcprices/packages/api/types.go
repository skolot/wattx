package api

type PriceRequest struct {
	Symbol  string `url:"symbol"`
	Convert string `url:"convert"`
}

type Status struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type Quote struct {
	Price float32 `json:"price"`
}

type Data struct {
	Quote map[string]Quote `json:"quote"`
}

type PriceResponse struct {
	Status Status          `json:"status"`
	Data   map[string]Data `json:"data"`
}

type PriceData map[string]map[string]float32
