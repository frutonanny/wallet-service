//go:generate mockgen --source=service.go --destination=mock/service.go
package get_balance

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
	internalErrors "github.com/frutonanny/wallet-service/internal/services/errors"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type Repository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
	GetBalance(ctx context.Context, walletID int64) (int64, error)
}

// RepoBuilder умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type RepoBuilder interface {
	NewRepository(db postgres.Database) Repository
}

type Service struct {
	logger      logger
	db          *sql.DB
	repoBuilder RepoBuilder
}

func New(logger logger, db *sql.DB) *Service {
	return &Service{
		logger:      logger,
		db:          db,
		repoBuilder: &builderImpl{},
	}
}

func (s *Service) WithBuilder(builder RepoBuilder) *Service {
	s.repoBuilder = builder
	return s
}

// GetBalance - отдает баланс пользователя.
// Алгоритм действий:
// - проверяем есть ли кошелек у пользователя, если нет, то возвращаем ошибку - ErrWalletNotFound;
// - отдаем баланс пользователя.
func (s *Service) GetBalance(ctx context.Context, userID int64) (int64, error) {
	repo := s.repoBuilder.NewRepository(s.db)

	// Проверяем есть ли кошелек у пользователя.
	walletID, err := repo.ExistWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoWalletNotFound) {
			s.logger.Error(fmt.Sprintf("for user %d wallet not found", userID))
			return 0, internalErrors.ErrWalletNotFound
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

	s.logger.Info(fmt.Sprintf("get balance of successfully for wallet: %d", walletID))

	return balance, nil
}
