package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	codes "github.com/frutonanny/wallet-service/pkg/responce_codes"
)

func (h *Handlers) PostAdd(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.AddRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.AddResponse{
			Error: &v1.Error{
				Code:    codes.InternalError,
				Message: "internal server error",
			},
		})
	}

	balance, err := h.addService.Add(ctx, req.UserID, req.Cash)
	if err != nil {
		return eCtx.JSON(http.StatusOK, v1.AddResponse{
			Error: &v1.Error{
				Code:    codes.InternalError,
				Message: "internal server error",
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.AddResponse{
		Data: &v1.AddData{
			Balance: balance,
		},
	})
}
