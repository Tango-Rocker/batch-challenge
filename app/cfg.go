package app

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	SourcePath string `env:"DATA_PATH"`
	SchemaPath string `env:"SCHEMA_PATH"`
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return cfg, err
	}

	return cfg, nil
}
