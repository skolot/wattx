package httpsrv

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"wattx/cmcprices/packages/api"
	"wattx/cmcprices/packages/config"
)

const (
	querySymbol string = "symbol"
	querConvert string = "convert"
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

	log.Printf("req.URL.Query(): %+v\n", req.URL.Query())

	symbol := q.Get(querySymbol)
	if symbol == "" {
		reportErr(w, errors.New("symbol is required"))
		return
	}

	convert := strings.Split(q.Get(querConvert), ",")
	if len(convert) == 0 {
		convert = []string{conf.API.Currency}
	}

	ops := api.PriceRequest{
		Symbol:  strings.Split(symbol, ","),
		Convert: convert,
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
	log.Println("Error: Price: " + err.Error())

	http.Error(w, err.Error(), http.StatusBadRequest)
}
