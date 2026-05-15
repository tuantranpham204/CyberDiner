package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/app"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/logger"
)

func main() {
	envFile := app.LoadDotenv()

	cfg, err := app.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.Init(cfg.Server.Mode); err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	if envFile != "" {
		logger.L().Infow("dotenv_loaded", "path", envFile)
	} else {
		logger.L().Info("dotenv_not_found_using_os_env")
	}

	a, err := app.New(cfg)
	if err != nil {
		logger.L().Fatalw("failed to initialize app", "error", err)
	}
	SetupRoutes(a)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:           a.Router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.L().Infow("server_starting", "port", cfg.Server.Port, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.L().Fatalw("server_error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.L().Info("server_shutting_down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.L().Errorw("server_shutdown_error", "error", err)
	}
}
