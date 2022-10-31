package main

import (
	"context"
	"errors"
	"fmt"
	logger2 "github.com/frutonanny/wallet-service/internal/logger"
	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/services/add"
	"log"
	"net"
	"os/signal"
	"syscall"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	serverGen "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	"github.com/frutonanny/wallet-service/internal/services/get_balance"
)

const path = "config/config.json"

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	config := serviceConfig.Must(path)

	logger := logger2.Must()

	db := postgres.MustConnect(config.DB.DSN)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(fmt.Sprintf("close db error: %s", err))
		}
	}()

	postgres.MustMigrate(db)

	// Address.
	addr := net.JoinHostPort(config.Service.Host, config.Service.Port)

	// Swagger.
	swagger, err := serverGen.GetSwagger()
	if err != nil {
		return fmt.Errorf("get swagger: %v", err)
	}

	// Services.
	getBalanceService := get_balance.New(logger, db)
	addService := add.New(logger, db)

	srv, err := initServer(addr, swagger, getBalanceService, addService)
	if err != nil {
		return fmt.Errorf("init server: %v", err)
	}

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
