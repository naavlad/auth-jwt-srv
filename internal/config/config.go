package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
}

type DatabaseConfig struct {
	URL string `env:"DATABASE_URL" env-required:"true"`
}

type JWTConfig struct {
	Secret               string        `env:"JWT_SECRET" env-required:"true"`
	AccessTokenDuration  time.Duration `env:"JWT_ACCESS_TOKEN_DURATION" env-default:"15m"`
	RefreshTokenDuration time.Duration `env:"JWT_REFRESH_TOKEN_DURATION" env-default:"168h"` // 7 days
}

type ServerConfig struct {
	Port string `env:"SERVER_PORT" env-default:"8080"`
}

func Load() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &cfg, nil
}
