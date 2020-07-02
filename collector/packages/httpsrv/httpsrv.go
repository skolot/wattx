package httpsrv

import (
	"log"
	"net/http"
	"strconv"

	"wattx/collector/packages/api"
	"wattx/collector/packages/config"
)

const (
	queryFormat string = "format"
	queryLimit  string = "limit"
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

	format := q.Get(queryFormat)
	if format == "" {
		format = defaultFormat
	}

	limit, err := safeAtoI(q.Get(queryLimit), conf.API.Limit)
	if err != nil {
		reportErr(w, err)
		return
	}

	data, err := api.RequestAllData(limit, conf)
	if err != nil {
		reportErr(w, err)
		return
	}

	if err = writeFormattedData(w, format, data); err != nil {
		reportErr(w, err)
		return
	}
}

func safeAtoI(a string, def int) (int, error) {
	if a == "" {
		return def, nil
	}

	return strconv.Atoi(a)
}

func reportErr(w http.ResponseWriter, err error) {
	log.Println("Error: Collector: " + err.Error())

	http.Error(w, err.Error(), http.StatusBadRequest)
}
