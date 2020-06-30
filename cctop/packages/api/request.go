package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"

	"wattx/top/packages/config"
)

const (
	apiAuthHeader string = "authorization"
	respError     string = "Error"
)

func Request(opts TopRequest, conf config.Config) ([]CoinInfo, error) {

	query, err := toQuery(opts)
	if err != nil {
		return nil, err
	}

	url := conf.API.URL + "?" + query
	data, err := doGetReq(url, conf.API.TimeoutDuration)
	if err != nil {
		return nil, err
	}

	return parseResp(data)
}

func toQuery(opts interface{}) (string, error) {
	q, err := query.Values(opts)
	if err != nil {
		return "", err
	}

	return q.Encode(), nil
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

func parseResp(data []byte) ([]CoinInfo, error) {

	log.Printf("data: %s\n", data)

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
		d.CoinInfo.Rank = i + 1
		ranked = append(ranked, d.CoinInfo)
	}

	return ranked, nil
}

func hasError(resp TopResponse) error {
	if resp.Response != respError {
		return nil
	}

	return errors.New(resp.Message)
}
