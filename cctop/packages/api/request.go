package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"wattx/cctop/packages/config"
)

const (
	apiAuthHeader string = "authorization"
	respError     string = "Error"

	queryLimit string = "limit"
	queryPage  string = "page"
	queryTSYM  string = "tsym"
)

func Request(opts TopRequest, conf config.Config) ([]CoinInfo, error) {
	url := conf.API.URL + "?" + toQuery(opts)
	data, err := doGetReq(url, conf.API.TimeoutDuration)
	if err != nil {
		return nil, err
	}

	rankShift := opts.Page * opts.Limit

	return parseResp(data, rankShift)
}

func toQuery(opts TopRequest) string {
	query := url.Values{}
	query.Add(queryLimit, strconv.Itoa(opts.Limit))
	query.Add(queryPage, strconv.Itoa(opts.Page))
	query.Add(queryTSYM, opts.TSYM)

	return query.Encode()
}

func doGetReq(url string, timeout time.Duration) ([]byte, error) {
	client := http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseResp(data []byte, rankShift int) ([]CoinInfo, error) {
	resp := TopResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if err := hasError(resp); err != nil {
		return nil, err
	}

	ciData := []Data{}
	if err := json.Unmarshal(resp.Data, &ciData); err != nil {
		return nil, err
	}

	ranked := []CoinInfo{}
	for i, d := range ciData {
		d.CoinInfo.Rank = i + 1 + rankShift
		ranked = append(ranked, d.CoinInfo)
	}

	log.Printf("ranked: %+v\n", ranked)

	return ranked, nil
}

func hasError(resp TopResponse) error {
	if resp.Response != respError {
		return nil
	}

	return errors.New(resp.Message)
}
