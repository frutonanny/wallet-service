package main

import (
	"github.com/getkin/kin-openapi/openapi3"

	server "github.com/frutonanny/wallet-service/internal/server/v1"
	"github.com/frutonanny/wallet-service/internal/server/v1/handlers"
	"github.com/frutonanny/wallet-service/internal/services/add"
	"github.com/frutonanny/wallet-service/internal/services/get_balance"
)

func initServer(
	addr string,
	swagger *openapi3.T,
	getBalanceService *get_balance.Service,
	addService *add.Service,
) (*server.Server, error) {
	h := handlers.NewHandlers(getBalanceService, addService)

	srv := server.New(
		addr,
		h,
		swagger,
	)

	return srv, nil
}
