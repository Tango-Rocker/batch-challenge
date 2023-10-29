package utils

import (
	"fmt"
	"github.com/caarlos0/env/v10"
)

func LoadEnvConfig[T any]() (*T, error) {
	cfg := new(T)
	if err := env.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
		return cfg, err
	}

	return cfg, nil
}
