package config

import (
	"os"

	"urm/internal/server"
)

type Config struct{}

func Load() *Config {
	return &Config{}
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
