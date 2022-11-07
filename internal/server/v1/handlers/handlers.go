package handlers

import (
	"context"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time"
	"time"

	"github.com/frutonanny/wallet-service/internal/services/get_transactions"
)

type getBalanceService interface {
	GetBalance(ctx context.Context, userID int64) (int64, error)
}
type addService interface {
	Add(ctx context.Context, walletID int64, cash int64) (int64, error)
}

type reserveService interface {
	Reserve(ctx context.Context, userID, serviceID, externalID, price int64) (int64, error)
}

type writeOffService interface {
	WriteOff(ctx context.Context, userID, serviceID, externalID, price int64) (int64, error)
}

type cancelService interface {
	Cancel(ctx context.Context, userID, orderID int64) (int64, error)
}

type getTransactions interface {
	GetTransactions(
		ctx context.Context,
		userID, limit, offset int64,
		sortBy get_transactions.SortBy,
		direction get_transactions.Direction,
	) ([]get_transactions.Transaction, error)
}
type getTransactionsByTime interface {
	GetTransactionsByTime(ctx context.Context, userID int64, start, end time.Time) ([]get_transactions_by_time.Transaction, error)
}

type getReport interface {
	GetReport(ctx context.Context, period string) (string, error)
}

type Handlers struct {
	getBalanceService     getBalanceService
	addService            addService
	reserveService        reserveService
	writeOffService       writeOffService
	cancelService         cancelService
	getTransactions       getTransactions
	getTransactionsByTime getTransactionsByTime
	getReport             getReport
}

func NewHandlers(
	getBalanceService getBalanceService,
	addService addService,
	reserveService reserveService,
	writeOffService writeOffService,
	cancelService cancelService,
	getTransactions getTransactions,
	getTransactionsByTime getTransactionsByTime,
	getReport getReport,
) *Handlers {
	return &Handlers{
		getBalanceService:     getBalanceService,
		addService:            addService,
		reserveService:        reserveService,
		writeOffService:       writeOffService,
		cancelService:         cancelService,
		getTransactions:       getTransactions,
		getTransactionsByTime: getTransactionsByTime,
		getReport:             getReport,
	}
}
