package httpsrv

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"wattx/cmcprices/packages/api"
	"wattx/cmcprices/packages/config"
)

const (
	symbolQuery   string = "symbol"
	currencyQuery string = "currency"
)

func Start(conf config.Config) error {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rootHandler(conf, w, req)
	})

	bind := conf.HTTPSrv.Host + ":" + strconv.Itoa(conf.HTTPSrv.Port)

	log.Println("Bind: http://" + bind)

	return http.ListenAndServe(bind, nil)
}

func rootHandler(conf config.Config, w http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()

	symbol := q.Get(symbolQuery)
	if symbol == "" {
		reportErr(w, errors.New("symbol is required"))
		return
	}

	currency := q.Get(currencyQuery)
	if currency == "" {
		currency = conf.API.Currency
	}

	ops := api.PriceRequest{
		Symbol:  symbol,
		Convert: currency,
	}

	data, err := api.Request(ops, conf)
	if err != nil {
		reportErr(w, err)
		return
	}

	if err = json.NewEncoder(w).Encode(data); err != nil {
		reportErr(w, err)
		return
	}
}

func reportErr(w http.ResponseWriter, err error) {
	log.Println("Error: Price Handler: " + err.Error())

	http.Error(w, err.Error(), http.StatusBadRequest)
}
