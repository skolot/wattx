package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"wattx/collector/packages/config"
)

const (
	apiAuthHeader string = "authorization"
	respError     string = "Error"

	queryLimit string = "limit"
	queryPage  string = "page"
	queryTSYM  string = "tsym"

	querySymbol  string = "symbol"
	queryConvert string = "convert"
)

func RequestAllData(limit int, conf config.Config) ([]Data, error) {
	data := []Data{}

	if limit <= 0 {
		limit = conf.API.Limit
	}

	requestSize := conf.API.RequestSize
	if requestSize > limit {
		requestSize = limit
	}

	pages := limit / requestSize

	for page := 0; page < pages; page++ {
		rrOpts := RankRequest{
			Limit: requestSize,
			Page:  page,
			TSYM:  conf.API.Currency,
		}

		rdata, err := RequestRank(rrOpts, conf)
		if err != nil {
			return nil, err
		}

		prOpts := PriceRequest{
			Symbol:  getSymbols(rdata),
			Convert: []string{conf.API.Currency},
		}

		pdata, err := RequestPrice(prOpts, conf)
		if err != nil {
			return nil, err
		}

		data = append(data, mergeData(rdata, pdata, conf.API.Currency)...)
	}

	return data, nil
}

func RequestPrice(opts PriceRequest, conf config.Config) (PriceResponse, error) {
	query := url.Values{}
	query.Add(querySymbol, strings.Join(opts.Symbol, ","))
	query.Add(queryConvert, strings.Join(opts.Convert, ","))

	url := conf.API.PriceURL + "?" + query.Encode()

	log.Printf("price request url: %s\n", url)

	data, err := Request(opts, url, conf)
	if err != nil {
		return nil, err
	}

	priceData, err := parsePriceResp(data)

	log.Printf("price response: %+v\n", priceData)

	return priceData, err
}

func RequestRank(opts RankRequest, conf config.Config) ([]RankResponse, error) {
	query := url.Values{}
	query.Add(queryLimit, strconv.Itoa(opts.Limit))
	query.Add(queryPage, strconv.Itoa(opts.Page))
	query.Add(queryTSYM, opts.TSYM)

	url := conf.API.TopURL + "?" + query.Encode()

	log.Printf("rank request url: %s\n", url)

	data, err := Request(opts, url, conf)
	if err != nil {
		return nil, err
	}

	rankData, err := parseRankResp(data)

	log.Printf("rank response: %+v\n", rankData)

	return rankData, err
}

func Request(opts interface{}, url string, conf config.Config) ([]byte, error) {
	data, statusCode, err := doGetReq(url, conf.API.TimeoutDuration)
	if err != nil {
		return nil, err
	}

	log.Printf("data: %s\n", data)

	if statusCode != http.StatusOK {
		if data != nil {
			return nil, parseErrorResp(data)
		}

		return nil, errors.New("requst failed with error code: " + strconv.Itoa(statusCode))
	}

	return data, nil
}

func doGetReq(url string, timeout time.Duration) ([]byte, int, error) {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	return data, resp.StatusCode, err
}

func parseRankResp(data []byte) ([]RankResponse, error) {
	resp := []RankResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func parsePriceResp(data []byte) (PriceResponse, error) {
	resp := PriceResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func parseErrorResp(data []byte) error {
	return errors.New(bytes.NewBuffer(data).String())
}

func getSymbols(rankData []RankResponse) []string {
	symbols := []string{}

	for _, rank := range rankData {
		symbols = append(symbols, rank.Name)
	}

	return symbols
}

func mergeData(rankData []RankResponse, priceData PriceResponse, currency string) []Data {
	merged := []Data{}

	for _, rank := range rankData {
		data := Data{}

		data.Rank = rank.Rank
		data.Price = priceData[rank.Name][currency]
		data.Name = rank.Name
		data.FullName = rank.FullName
		data.Currency = currency

		merged = append(merged, data)
	}

	return merged
}
