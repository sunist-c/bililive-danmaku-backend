package model

import (
	"flag"
)

type Config struct {
	BackendMode bool
	ServingPort uint
	AccessToken string
}

func NewConfigWithBackendMode(servingPort uint, accessToken string) *Config {
	return &Config{
		BackendMode: true,
		ServingPort: servingPort,
		AccessToken: accessToken,
	}
}

func NewConfigWithSingleMode() *Config {
	return &Config{
		BackendMode: false,
		ServingPort: 0,
		AccessToken: "",
	}
}

func init() {
	if config == nil {
		config = &Config{}
	}

	flag.BoolVar(&config.BackendMode, "b", true, "enable backend mode")
	flag.UintVar(&config.ServingPort, "p", 8080, "port to serving on")
	flag.StringVar(&config.AccessToken, "t", "", "access-token for authentication")

	flag.Parse()
}
