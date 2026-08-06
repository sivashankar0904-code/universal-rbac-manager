package config

import (
	"os"

	"urm/internal/server"
)

type Config struct {
	Env      string
	LogLevel string
}

func Load() *Config {
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
		Host: c.getEnv("URM_HTTP_HOST", "0.0.0.0"),
		Port: c.getEnv("URM_HTTP_PORT", "8080"),
	}
}
