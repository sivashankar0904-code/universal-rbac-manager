package main

import (
	"log/slog"
	"os"

	"urm/internal/config"
	"urm/internal/logger"
	"urm/internal/server"
)

func startAllServices(log *slog.Logger, cfg *config.Config) {
	srv := server.New(cfg.HTTP(), log)
	if err := srv.Start(); err != nil {
		log.Error("failed to start http server", "err", err)
		os.Exit(1)
	}
}

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env, cfg.LogLevel)

	startAllServices(log, cfg)
}
