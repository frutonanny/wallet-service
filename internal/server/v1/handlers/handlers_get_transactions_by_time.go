package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time"
	"github.com/frutonanny/wallet-service/pkg/errcodes"
)

func (h *Handlers) PostGetTransactionsByTime(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetTransactionsByTimeRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetTransactionsByTimeResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	txs, err := h.getTransactionsByTime.GetTransactionsByTime(ctx, req.UserID, req.Start, req.End)
	if err != nil {
		code := errcodes.InternalError
		msg := "internal server error"

		if errors.Is(err, servicesErrors.ErrWalletNotFound) {
			code = errcodes.WalletNotFound
			msg = "wallet not found"
		}

		return eCtx.JSON(http.StatusOK, v1.GetTransactionsByTimeResponse{
			Error: &v1.Error{
				Code:    code,
				Message: msg,
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetTransactionsByTimeResponse{
		Data: &v1.GetTransactionsByTimeData{
			Transactions: adaptTxsByTime(txs),
		},
	})
}

func adaptTxsByTime(txs []get_transactions_by_time.Transaction) []v1.Transaction {
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
