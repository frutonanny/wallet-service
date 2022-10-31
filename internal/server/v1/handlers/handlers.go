package handlers

import (
	"context"
	"github.com/labstack/echo/v4"
)

type getBalanceService interface {
	GetBalance(ctx context.Context, userID int64) (int64, error)
}
type addService interface {
	Add(ctx context.Context, walletID int64, cash int64) (int64, error)
}

type Handlers struct {
	getBalanceService getBalanceService
	addService        addService
}

func NewHandlers(getBalanceService getBalanceService, addService addService) *Handlers {
	return &Handlers{
		getBalanceService: getBalanceService,
		addService:        addService,
	}
}

func (h *Handlers) PostCancel(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) PostGetReport(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) PostGetTransactions(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) PostGetTransactionsByTime(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) PostReserve(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *Handlers) PostWriteOff(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
