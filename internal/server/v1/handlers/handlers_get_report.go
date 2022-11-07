package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v1 "github.com/frutonanny/wallet-service/internal/generated/server/v1"
	"github.com/frutonanny/wallet-service/pkg/errcodes"
)

func (h *Handlers) PostGetReport(eCtx echo.Context) error {
	ctx := eCtx.Request().Context()

	var req v1.GetReportRequest
	if err := eCtx.Bind(&req); err != nil {
		return eCtx.JSON(http.StatusOK, v1.GetReportResponse{
			Error: &v1.Error{
				Code:    errcodes.InternalError,
				Message: "internal server error",
			},
		})
	}

	url, err := h.getReport.GetReport(ctx, req.Period)
	if err != nil {
		code := errcodes.InternalError
		msg := "internal server error"

		return eCtx.JSON(http.StatusOK, v1.GetReportResponse{
			Error: &v1.Error{
				Code:    code,
				Message: msg,
			},
		})
	}

	return eCtx.JSON(http.StatusOK, v1.GetReportResponse{
		Data: &v1.GetReportData{
			Url: url,
		},
	})
}
