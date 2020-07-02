package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"wattx/cmcprices/packages/config"
)

const (
	apiAuthHeader string = "X-CMC_PRO_API_KEY"

	querySymbol  string = "symbol"
	queryConvert string = "convert"

	retriesCount int           = 5
	retryTimeout time.Duration = 100 * time.Millisecond

	invalidSymbolsErrorPrefix string = "Invalid values for \"symbol\": \""
	invalidSymbolsErrorSuffix string = "\""
)

func Request(opts PriceRequest, conf config.Config) (PriceData, error) {
	opts0 := opts

	for try := 0; try < retriesCount; try++ {
		url := conf.API.URL + "?" + toQuery(opts0)
		log.Printf("price request url: %s\n", url)

		data, err := doGetReq(url, conf)
		if err != nil {
			return nil, err
		}

		priceData, invalidSymbols, err := parseResp(data)
		if err != nil && len(invalidSymbols) > 0 {
			log.Println("retrying, got symbols error", err)
			opts0.Symbol = dropInvalidSymbols(opts0.Symbol, invalidSymbols)
			continue
		}

		if err != nil {
			return nil, err
		}

		log.Printf("price response: %+v\n", priceData)

		return priceData, nil
	}

	return nil, errors.New("Too many errors")
}

func toQuery(opts PriceRequest) string {
	query := url.Values{}
	query.Add(querySymbol, strings.Join(opts.Symbol, ","))
	query.Add(queryConvert, strings.Join(opts.Convert, ","))

	return query.Encode()
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

func parseResp(data []byte) (PriceData, map[string]bool, error) {
	resp := PriceResponse{}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, map[string]bool{}, err
	}

	if invalidSymbols, err := hasError(resp); err != nil {
		return nil, invalidSymbols, err
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

	return priceData, map[string]bool{}, nil
}

func hasError(resp PriceResponse) (map[string]bool, error) {
	if resp.Status.ErrorCode == 0 {
		return map[string]bool{}, nil
	}

	return getInvalidSymbols(resp), errors.New(resp.Status.ErrorMessage)
}

func getInvalidSymbols(resp PriceResponse) map[string]bool {
	if strings.HasPrefix(resp.Status.ErrorMessage, invalidSymbolsErrorPrefix) {
		invalidSymbols := strings.TrimSuffix(
			strings.TrimPrefix(resp.Status.ErrorMessage, invalidSymbolsErrorPrefix),
			invalidSymbolsErrorSuffix,
		)

		invalidSymbolsMap := map[string]bool{}
		for _, symbol := range strings.Split(invalidSymbols, ",") {
			invalidSymbolsMap[symbol] = true
		}

		return invalidSymbolsMap
	}

	return map[string]bool{}
}

func dropInvalidSymbols(symbols []string, invalidSymbols map[string]bool) []string {
	if len(invalidSymbols) == 0 {
		return symbols
	}

	filtered := []string{}

	for _, symbol := range symbols {
		if _, ok := invalidSymbols[symbol]; ok {
			continue
		}

		filtered = append(filtered, symbol)
	}

	return filtered
}
