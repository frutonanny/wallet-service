package handlers

import (
	"context"

	"github.com/labstack/echo/v4"

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

type Handlers struct {
	getBalanceService getBalanceService
	addService        addService
	reserveService    reserveService
	writeOffService   writeOffService
	cancelService     cancelService
	getTransactions   getTransactions
}

func NewHandlers(
	getBalanceService getBalanceService,
	addService addService,
	reserveService reserveService,
	writeOffService writeOffService,
	cancelService cancelService,
	getTransactions getTransactions,
) *Handlers {
	return &Handlers{
		getBalanceService: getBalanceService,
		addService:        addService,
		reserveService:    reserveService,
		writeOffService:   writeOffService,
		cancelService:     cancelService,
		getTransactions:   getTransactions,
	}
}

func (h *Handlers) PostGetTransactionsByTime(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
