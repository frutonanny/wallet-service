package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	server "github.com/frutonanny/wallet-service/internal/server/v1"
	"github.com/frutonanny/wallet-service/internal/server/v1/handlers"
	"github.com/frutonanny/wallet-service/internal/services/add"
	"github.com/frutonanny/wallet-service/internal/services/cancel"
	"github.com/frutonanny/wallet-service/internal/services/get_balance"
	"github.com/frutonanny/wallet-service/internal/services/get_report"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time"
	"github.com/frutonanny/wallet-service/internal/services/reserve"
	write_off "github.com/frutonanny/wallet-service/internal/services/write-off"
)

func initServer(
	addr string,
	swagger *openapi3.T,
	getBalanceService *get_balance.Service,
	addService *add.Service,
	reserveService *reserve.Service,
	writeOffService *write_off.Service,
	cancelService *cancel.Service,
	getTransactions *get_transactions.Service,
	getTransactionsByTime *get_transactions_by_time.Service,
	getReport *get_report.Service,
) (*server.Server, error) {
	h := handlers.NewHandlers(
		getBalanceService,
		addService,
		reserveService,
		writeOffService,
		cancelService,
		getTransactions,
		getTransactionsByTime,
		getReport,
	)

	srv := server.New(
		addr,
		h,
		swagger,
	)

	return srv, nil
}
