package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/caarlos0/env/v11"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
)

type Config struct {
	Env         Environment `env:"APP_ENV" envDefault:"development"`
	HTTPPort    int         `env:"HTTP_PORT"  envDefault:"8080"`
	LogLevel    slog.Level  `env:"LOG_LEVEL"  envDefault:"info"`
	DatabaseURL string      `env:"DATABASE_URL"`
}

func Load() (Config, error) {
	cfg, err := env.ParseAs[Config]()

	if err != nil {
		return Config{}, fmt.Errorf("parse environment variables: %w", err)
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("parse environment variables: %w", err)
	}

	return cfg, nil
}

func (c Config) IsProduction() bool {
	return c.Env == EnvProduction
}
func (c Config) HTTPAddr() string {
	return ":" + strconv.Itoa(c.HTTPPort)
}

func (c Config) validate() error {
	var errs []error

	switch c.Env {
	case EnvDevelopment, EnvProduction:
	default:
		errs = append(errs, fmt.Errorf("APP_ENV must be between %q or %q, got %q", EnvDevelopment, EnvProduction, c.Env))
	}

	if c.HTTPPort < 1 || c.HTTPPort > 65535 {
		errs = append(errs, fmt.Errorf("HTTP_PORT must be between 1 and 65535, got %d", c.HTTPPort))
	}

	if c.DatabaseURL == "" {
		errs = append(errs, errors.New("DATABASE_URL cannot be empty"))
	}

	return errors.Join(errs...)
}
