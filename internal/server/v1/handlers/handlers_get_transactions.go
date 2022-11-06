package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions"
	"github.com/frutonanny/wallet-service/pkg/errcodes"
)

func (h *Handlers) PostGetTransactions(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetTransactionsRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetTransactionsResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	txs, err := h.getTransactions.
		GetTransactions(
			ctx,
			req.UserID,
			req.Limit,
			req.Offset,
			adaptSortBy(req.SortBy),
			adaptDirection(req.Direction),
		)

	if err != nil {
		code := errcodes.InternalError
		msg := "internal server error"

		if errors.Is(err, servicesErrors.ErrWalletNotFound) {
			code = errcodes.WalletNotFound
			msg = "wallet not found"
		}

		return eCtx.JSON(http.StatusOK, v1.GetTransactionsResponse{
			Error: &v1.Error{
				Code:    code,
				Message: msg,
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetTransactionsResponse{
		Data: &v1.GetTransactionsData{
			Transactions: adaptTxs(txs),
		},
	})
}

func adaptTxs(txs []get_transactions.Transaction) []v1.Transaction {
	result := make([]v1.Transaction, 0, len(txs))

	for i := range txs {
		tx := v1.Transaction{
			Amount:      txs[i].Amount,
			CreatedAt:   txs[i].CreatedAt,
			Description: txs[i].Description,
		}

		result = append(result, tx)
	}

	return result
}

func adaptSortBy(sortBy v1.GetTransactionsRequestSortBy) get_transactions.SortBy {
	switch sortBy {
	case v1.Amount:
		return get_transactions.Amount
	default:
		return get_transactions.Date
	}
}

func adaptDirection(direction v1.GetTransactionsRequestDirection) get_transactions.Direction {
	switch direction {
	case v1.Asc:
		return get_transactions.Asc
	default:
		return get_transactions.Desc
	}
}
