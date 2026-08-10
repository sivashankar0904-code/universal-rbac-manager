package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"urm/internal/apm"
	"urm/internal/config"
	"urm/internal/helper"
	"urm/internal/logger"
	"urm/internal/repo"
	"urm/internal/server"
	"urm/internal/service"
)

func startAllServices(log *slog.Logger, cfg *config.Config, otelCfg apm.Config) {
	ctx := context.Background()

	shutdown, err := apm.Init(ctx, otelCfg)
	if err != nil {
		// logger isn't built yet at this point, hence the bare stderr write
		fmt.Fprintln(os.Stderr, "failed to init apm:", err)
		os.Exit(1)
	}
	defer shutdown(ctx)

	pool, err := repo.NewPool(ctx, *cfg.Database(), log)
	if err != nil {
		log.Error("failed to connect to database", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	//helper service
	helperSvc := helper.NewService()
	//repo service contain all db operations
	repoSvc := repo.NewService(pool, helperSvc)
	//services is grouped service with different interconnected packages
	serviceSvc := service.NewService(repoSvc, helperSvc)

	handler := server.NewHandler(
		log,
		serviceSvc,
	)
	log.Info("services initialized")

	srv := server.New(cfg.HTTP(), log, handler, cfg.APM().ServiceName)
	if err := srv.Start(); err != nil {
		log.Error("failed to start http server", "err", err)
		os.Exit(1)
	}
}

func main() {
	cfg := config.Load()
	otelCfg := *cfg.APM()

	log := logger.New(cfg.Env, cfg.LogLevel, otelCfg.ServiceName)

	startAllServices(log, cfg, otelCfg)
}
