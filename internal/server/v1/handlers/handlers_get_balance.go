package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/pkg/errcodes"
)

func (h *Handlers) PostGetBalance(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetBalanceRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	balance, err := h.getBalanceService.GetBalance(ctx, req.UserID)
	if err != nil {
		code := errcodes.InternalError
		msg := "internal server error"

		if errors.Is(err, servicesErrors.ErrWalletNotFound) {
			code = errcodes.WalletNotFound
			msg = "wallet not found"
		}

		return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
			Error: &v1.Error{
				Code:    code,
				Message: msg,
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
		Data: &v1.GetBalanceData{
			Balance: balance,
		},
	})
}
