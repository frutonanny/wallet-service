//go:generate mockgen --source=service.go --destination=mock/service.go
package get_transactions_by_time

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
	"github.com/frutonanny/wallet-service/internal/repositories/transaction"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type WalletRepository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
}

type TransactionRepository interface {
	GetTransactionsByTime(
		ctx context.Context,
		walletID int64,
		timeStart, timeEnd time.Time,
	) ([]transaction.Transaction, error)
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewWalletRepository(db postgres.Database) WalletRepository
	NewTransactionRepository(db postgres.Database) TransactionRepository
}

type Service struct {
	db     *sql.DB
	logger logger
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

// GetTransactionsByTime - отдает список транзакций пользователя, отсортированный по переданному параметру.
// - проверяем есть ли кошелек у пользователя, если нет, то отдаем ошибку ErrWalletNotFound.
// - отдает список транзакций.
func (s *Service) GetTransactionsByTime(ctx context.Context, userID int64, start, end time.Time) ([]Transaction, error) {
	walletRepo := s.deps.NewWalletRepository(s.db)

	// Проверяем есть ли кошелек у пользователя.
	walletID, err := walletRepo.ExistWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoWalletNotFound) {
			s.logger.Error(fmt.Sprintf("for user %d wallet not found", userID))
			return nil, servicesErrors.ErrWalletNotFound
		}

		s.logger.Error(fmt.Sprintf("exist wallet: %s", err))
		return nil, fmt.Errorf("exist wallet: %w", err)
	}

	txsRepo := s.deps.NewTransactionRepository(s.db)

	// Отдаем список транзакций, отсортированный по переданному параметру.
	txs, err := txsRepo.GetTransactionsByTime(ctx, walletID, start, end)
	if err != nil {
		s.logger.Error(fmt.Sprintf("get transactions: %s", err))
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	result, err := adaptTxs(txs)
	if err != nil {
		return nil, fmt.Errorf("adapt txs: %v", err)
	}

	return result, nil
}
