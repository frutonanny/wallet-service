//go:generate mockgen --source=service.go --destination=mock/service.go
package get_balance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type Repository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
	GetBalance(ctx context.Context, walletID int64) (int64, error)
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewRepository(db postgres.Database) Repository
}

type Service struct {
	logger logger
	db     *sql.DB
	deps   dependencies
}

func New(logger logger, db *sql.DB) *Service {
	return &Service{
		logger: logger,
		db:     db,
		deps:   &dependenciesImpl{},
	}
}

func (s *Service) WithDependencies(deps dependencies) *Service {
	s.deps = deps
	return s
}

// GetBalance - отдает баланс пользователя.
// - проверяем есть ли кошелек у пользователя, если нет, то возвращаем ошибку - ErrWalletNotFound;
// - отдаем баланс пользователя.
func (s *Service) GetBalance(ctx context.Context, userID int64) (int64, error) {
	repo := s.deps.NewRepository(s.db)

	// Проверяем есть ли кошелек у пользователя.
	walletID, err := repo.ExistWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoWalletNotFound) {
			s.logger.Error(fmt.Sprintf("for user %d wallet not found", userID))
			return 0, servicesErrors.ErrWalletNotFound
		}

		s.logger.Error(fmt.Sprintf("exist wallet: %s", err))
		return 0, fmt.Errorf("exist wallet: %w", err)
	}

	// Получаем текущий баланс пользователя.
	balance, err := repo.GetBalance(ctx, walletID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("get balance: %s", err))
		return 0, fmt.Errorf("get balance: %w", err)
	}

	return balance, nil
}
