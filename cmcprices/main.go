package main

import (
	"log"
	"os"

	"wattx/cmcprices/packages/config"
	"wattx/cmcprices/packages/httpsrv"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Println("Error: config: ", err)
		os.Exit(1)
	}

	httpsrv.Start(conf)
}
