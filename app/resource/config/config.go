package config

import (
	"github.com/caarlos0/env/v9"
	envs "github.com/simpleAI/service-video-maker/app/resource/constants/env"
)

type Config struct {
	Environment envs.Environment `env:"ENVIRONMENT" envDefault:"DEVELOPMENT"`
	Port        string           `env:"PORT" envDefault:"9000"`
}

func New() (*Config, error) {
	cfg := new(Config)

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
