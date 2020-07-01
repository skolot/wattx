package config

import (
	"time"

	"github.com/jinzhu/configor"
)

type API struct {
	URL             string
	Timeout         string
	TimeoutDuration time.Duration
	Key             string
	Currency        string
}

type HTTPSrv struct {
	Host string
	Port int
}

type Config struct {
	API     API
	HTTPSrv HTTPSrv
}

const (
	configFile string = "config/config.json"
)

func Read() (Config, error) {
	conf := Config{}

	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&conf, configFile)
	if err != nil {
		return Config{}, err
	}

	conf.API.TimeoutDuration, err = time.ParseDuration(conf.API.Timeout)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
