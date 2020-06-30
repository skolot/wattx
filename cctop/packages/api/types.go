package api

import (
	"encoding/json"
)

type TopRequest struct {
	Limit int    `url:"limit"`
	Page  int    `url:"page"`
	TSYM  string `url:"tsym"`
}

type ErrorCode struct {
	Response string   `json:"Response"`
	Message  string   `json:"Message"`
	Data     struct{} `json:"Data"`
}

type CoinInfo struct {
	Rank     int    `json:"Rank"`
	Id       string `json:"Id"`
	Name     string `json:"Name"`
	FullName string `json:"FullName"`
}

type Data struct {
	CoinInfo CoinInfo `json:"CoinInfo"`
}

type TopResponse struct {
	ErrorCode
	Data json.RawMessage `json:"Data"`
}
