package v1

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	mdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"

	"github.com/frutonanny/wallet-service/internal/generated/server/v1"
)

const (
	readHeaderTimeout = 10 * time.Second
	defaultTimeout    = 3 * time.Second
)

type Server struct {
	srv *http.Server
}

func New(
	addr string,
	handlers v1.ServerInterface,
	swagger *openapi3.T,
) *Server {
	e := echo.New()

	group := e.Group("v1", mdlwr.OapiRequestValidator(swagger))
	v1.RegisterHandlers(group, handlers)

	return &Server{
		srv: &http.Server{
			Addr:              addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		_ = s.srv.Shutdown(ctx)
	}()

	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}