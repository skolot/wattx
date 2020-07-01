package httpsrv

import (
	"log"
	"net/http"
	"strconv"

	"wattx/collector/packages/api"
	"wattx/collector/packages/config"
)

const (
	formatQuery string = "format"
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

	format := q.Get(formatQuery)
	if format == "" {
		format = defaultFormat
	}

	data, err := api.RequestAllData(conf)
	if err != nil {
		reportErr(w, err)
		return
	}

	if err = writeFormattedData(w, format, data); err != nil {
		reportErr(w, err)
		return
	}
}

func reportErr(w http.ResponseWriter, err error) {
	log.Println("Error: Collector: " + err.Error())

	http.Error(w, err.Error(), http.StatusBadRequest)
}
