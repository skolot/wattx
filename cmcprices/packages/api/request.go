package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"

	"wattx/cmcprices/packages/config"
)

const (
	apiAuthHeader string = "X-CMC_PRO_API_KEY"
)

func Request(opts PriceRequest, conf config.Config) (PriceData, error) {
	query, err := toQuery(opts)
	if err != nil {
		return nil, err
	}

	url := conf.API.URL + "?" + query
	data, err := doGetReq(url, conf)
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

func doGetReq(url string, conf config.Config) ([]byte, error) {
	client := http.Client{
		Timeout: conf.API.TimeoutDuration,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accepts", "application/json")
	req.Header.Add(apiAuthHeader, conf.API.Key)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func parseResp(data []byte) (PriceData, error) {
	resp := PriceResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	if err := hasError(resp); err != nil {
		return nil, err
	}

	priceData := PriceData{}

	for name, data := range resp.Data {
		coinPrices, ok := priceData[name]
		if !ok {
			coinPrices = map[string]float32{}
		}

		for currency, quote := range data.Quote {
			coinPrices[currency] = quote.Price
		}

		priceData[name] = coinPrices
	}

	return priceData, nil
}

func hasError(resp PriceResponse) error {
	if resp.Status.ErrorCode == 0 {
		return nil
	}

	return errors.New(resp.Status.ErrorMessage)
}
