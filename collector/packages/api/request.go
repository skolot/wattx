package api

import (
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

func RequestAllData(conf config.Config) ([]Data, error) {
	data := []Data{}

	limit := conf.API.Limit
	if limit > conf.API.Top {
		limit = conf.API.Top
	}

	pages := conf.API.Top / limit

	for page := 0; page < pages; page++ {
		rrOpts := RankRequest{
			Limit: limit,
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

	log.Printf("conf.API.PriceURL ? query.Encode(): %s\n", conf.API.PriceURL+"?"+query.Encode())

	data, err := Request(opts, conf.API.PriceURL+"?"+query.Encode(), conf)
	if err != nil {
		return nil, err
	}

	return parsePriceResp(data)
}

func RequestRank(opts RankRequest, conf config.Config) ([]RankResponse, error) {
	query := url.Values{}
	query.Add(queryLimit, strconv.Itoa(opts.Limit))
	query.Add(queryPage, strconv.Itoa(opts.Page))
	query.Add(queryTSYM, opts.TSYM)

	data, err := Request(opts, conf.API.TopURL+"?"+query.Encode(), conf)
	if err != nil {
		return nil, err
	}

	return parseRankResp(data)
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
	resp := ErrorResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return err
	}

	return errors.New(resp.Error)
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
