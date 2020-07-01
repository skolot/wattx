package main

import (
	"log"

	"wattx/cmcprices/packages/config"
	"wattx/cmcprices/packages/httpsrv"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Println("Error: config: ", err)
	}

	httpsrv.Start(conf)
}
