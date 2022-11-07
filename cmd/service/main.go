package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os/signal"
	"syscall"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	serverGen "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	logger2 "github.com/frutonanny/wallet-service/internal/logger"
	"github.com/frutonanny/wallet-service/internal/minio"
	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/services/add"
	cancelSev "github.com/frutonanny/wallet-service/internal/services/cancel"
	"github.com/frutonanny/wallet-service/internal/services/get_balance"
	"github.com/frutonanny/wallet-service/internal/services/get_report"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time"
	"github.com/frutonanny/wallet-service/internal/services/reserve"
	write_off "github.com/frutonanny/wallet-service/internal/services/write-off"
)

var configFile string

func init() {
	flag.StringVar(&configFile,
		"config",
		"config/config.local.json",
		"Path to configuration file")
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	flag.Parse()

	f := flag.Lookup(serviceConfig.Arg)
	if f == nil {
		return errors.New("config arg must be set")
	}

	config := serviceConfig.Must(f.Value.String())

	logger := logger2.Must()

	// Postgres.
	db := postgres.MustConnect(config.DB.DSN)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(fmt.Sprintf("close db error: %s", err))
		}
	}()

	postgres.MustMigrate(db)

	// Minio.
	minioClient := minio.Must(
		config.Minio.Endpoint,
		config.Minio.AccessKeyID,
		config.Minio.SecretAccessKey,
	)

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
	reserveService := reserve.New(logger, db)
	writeOffService := write_off.New(logger, db)
	cancelService := cancelSev.New(logger, db)
	getTransactions := get_transactions.New(logger, db)
	getTransactionsByTime := get_transactions_by_time.New(logger, db)
	getReport := get_report.New(logger, db, minioClient, config.Minio.PublicEndpoint)

	srv, err := initServer(
		addr,
		swagger,
		getBalanceService,
		addService,
		reserveService,
		writeOffService,
		cancelService,
		getTransactions,
		getTransactionsByTime,
		getReport,
	)

	if err != nil {
		return fmt.Errorf("init server: %v", err)
	}

	if err := srv.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("run server: %v", err)
	}

	return nil
}
