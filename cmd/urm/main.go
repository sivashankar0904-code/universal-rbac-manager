package main

import (
	"context"
	"log/slog"
	"os"

	"urm/internal/config"
	"urm/internal/helper"
	"urm/internal/logger"
	"urm/internal/repo"
	"urm/internal/server"
	"urm/internal/service"
)

func startAllServices(log *slog.Logger, cfg *config.Config) {
	ctx := context.Background()

	pool, err := repo.NewPool(ctx, *cfg.Database(), log)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	helperSvc := helper.NewService()
	repoSvc := repo.NewService(pool, helperSvc)
	serviceSvc := service.NewService(repoSvc, helperSvc)

	handler := server.NewHandler(
		log,
		serviceSvc,
	)
	log.Info("services initialized")

	srv := server.New(cfg.HTTP(), log, handler)
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
