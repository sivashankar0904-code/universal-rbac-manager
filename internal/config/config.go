package config

import (
	"os"

	"github.com/joho/godotenv"

	"urm/internal/apm"
	"urm/internal/repo"
	"urm/internal/server"
)

type Config struct {
	Env      string
	LogLevel string
}

func Load() *Config {
	_ = godotenv.Load()

	c := &Config{}
	c.Env = c.getEnv("URM_ENV", "local")
	c.LogLevel = c.getEnv("URM_LOG_LEVEL", "info")
	return c
}

func (c *Config) getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c *Config) HTTP() *server.Config {
	return &server.Config{
		Host: c.getEnv("URM_HTTP_HOST", ""),
		Port: c.getEnv("URM_HTTP_PORT", ""),
	}
}

func (c *Config) Database() *repo.Config {
	return &repo.Config{
		Host:     c.getEnv("URM_DB_HOST", ""),
		Port:     c.getEnv("URM_DB_PORT", ""),
		User:     c.getEnv("URM_DB_USER", ""),
		Password: c.getEnv("URM_DB_PASSWORD", ""),
		Name:     c.getEnv("URM_DB_NAME", ""),
		SSLMode:  c.getEnv("URM_DB_SSLMODE", "disable"),
	}
}

func (c *Config) APM() *apm.Config {
	return &apm.Config{
		Enabled:           c.getEnv("URM_OTEL_ENABLED", "false") == "true",
		CollectorEndpoint: c.getEnv("URM_OTEL_COLLECTOR_ENDPOINT", "localhost:4317"),
		ServiceName:       c.getEnv("URM_OTEL_SERVICE_NAME", "urm"),
		ServiceVersion:    c.getEnv("URM_OTEL_SERVICE_VERSION", "dev"),
		Environment:       c.Env,
	}
}
