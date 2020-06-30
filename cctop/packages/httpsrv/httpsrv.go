package httpsrv

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"wattx/cctop/packages/api"
	"wattx/cctop/packages/config"
)

const (
	limitQuery string = "limit"
	pageQuery  string = "page"
	tsymQuery  string = "tsym"
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

	limit, err := safeAtoI(q.Get(limitQuery), conf.API.Limit)
	if err != nil {
		reportErr(w, err)
		return
	}

	page, err := safeAtoI(q.Get(pageQuery), 0)
	if err != nil {
		reportErr(w, err)
		return
	}

	tsym := q.Get(tsymQuery)
	if tsym == "" {
		tsym = conf.API.TSYM
	}

	ops := api.TopRequest{
		Limit: limit,
		Page:  page,
		TSYM:  tsym,
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

func safeAtoI(a string, def int) (int, error) {
	if a == "" {
		return def, nil
	}

	return strconv.Atoi(a)
}

func reportErr(w http.ResponseWriter, err error) {
	log.Println("Error: Top Handler: " + err.Error())

	http.Error(w, err.Error(), http.StatusBadRequest)
}
