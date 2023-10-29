package bussines

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	BufferSize   int `env:"BUFFER_SIZE" envDefault:"1000"`
	FlushTimeout int `env:"FLUSH_TIMEOUT_MS" envDefault:"1000"`
}

func LoadEnvConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return cfg, err
	}

	return cfg, nil
}
