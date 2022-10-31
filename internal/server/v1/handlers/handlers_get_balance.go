package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	codes "github.com/frutonanny/wallet-service/pkg/responce_codes"
)

func (h *Handlers) PostGetBalance(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetBalanceRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
			Error: &v1.Error{
				Code:    codes.InternalError,
				Message: "internal server error",
			},
		})
	}

	balance, err := h.getBalanceService.GetBalance(ctx, req.UserID)
	if err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
			Error: &v1.Error{
				Code:    codes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetBalanceResponse{
		Data: &v1.GetBalanceData{
			Balance: balance,
		},
	})
}
