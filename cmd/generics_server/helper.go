package main

import (
	"boilerplate-demo/src/interface/fiber_server"
	"boilerplate-demo/src/repository/player_repository"
	"boilerplate-demo/src/use_case"
	"boilerplate-demo/src/use_case/repository"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	fconfig "boilerplate-demo/src/interface/fiber_server/config"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logger *zap.Logger

func initEnvironment() config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Failed loading .env file: %s", err)
	}

	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parse env to struct: %v", err)
	}

	return cfg
}

func initLogger(cfg config) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)
	if cfg.Debuglog {
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	config.Level = logLevel

	lg, err := config.Build()
	if err != nil {
		log.Fatalf("Error build logger: %s\n", err)
	}
	defer lg.Sync()

	zap.ReplaceGlobals(lg)
	logger = zap.L().Named("demo")
	logger.Info("Logger initialized")
}

func initRepositories(cfg config) (playerRepository repository.PlayerRepository) {
	playerRepository = player_repository.NewGormPostgres(setupGorm(postgres.Open(cfg.Services.PostgresqlUri)), logger)

	return
}

func initInterfaces(cfg config, uc *use_case.UseCase) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := uc.HealthCheck(ctx); err != nil {
		logger.Fatal("Error health check fail before start", zap.Error(err))
	}

	wg := new(sync.WaitGroup)
	serv := fiber_server.New(uc, &fconfig.ServerConfig{
		AppVersion:    cfg.AppVersion,
		ListenAddress: fmt.Sprintf(":%d", cfg.Port),
		RequestLog:    true,
	})
	logger.Info("Fiber server initialized")

	serv.Start(wg)
	logger.Info("Fiber server started")

	wg.Wait()
	logger.Info("Application stopped")
}

func setupGorm(d gorm.Dialector) *gorm.DB {
	db, err := gorm.Open(d, &gorm.Config{})
	if err != nil {
		logger.Fatal("Error open postgres url", zap.Error(err))
	}
	return db
}
