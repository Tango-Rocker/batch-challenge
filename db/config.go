package db

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	Host string `env:"DB_HOST"`
	Port string `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Name string `env:"DB_NAME"`
}

func LoadEnvConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return cfg, err
	}

	return cfg, nil
}
