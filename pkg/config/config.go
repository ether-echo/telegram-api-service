package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
	//Token    string `env:"TOKEN" envDefault:"7731167021:AAHb5ofDkyroDH0LUYznmRVeIycJCYe0c2M"`
	Token string `env:"TOKEN" envDefault:"7746162400:AAF09997wm2OCcnvVcWK-V5q9fYzmrKWeXA"` //Test
}

func ReadConfig() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return &config, err
}
