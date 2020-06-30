package main

import (
	"log"

	"wattx/top/packages/config"
	"wattx/top/packages/httpsrv"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Println("Error: config")
	}

	httpsrv.Start(conf)
}
