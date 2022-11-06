//go:generate mockgen --source=service.go --destination=mock/service.go
package add

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type WalletRepository interface {
	CreateIfNotExist(ctx context.Context, userID int64) (int64, error)
	Add(ctx context.Context, walletID int64, cash int64) (int64, error)
}

type TransactionRepository interface {
	AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error
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

// Add - начисляет переданную сумму на счет пользователя.
// - проверяем есть ли кошелек у пользователя, если нет, то создаем;
// - зачисляем переданную сумму на кошелек пользователя;
// - добавляем транзакцию о внесенных средствах;
// - в ответ отдаем текущий баланс пользователя в копейках с учетом пополнения.
func (s *Service) Add(ctx context.Context, userID, cash int64) (int64, error) {
	// Стартуем транзакцию.
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.logger.Error(fmt.Sprintf("begin tx: %s", err))
		return 0, fmt.Errorf("begin tx: %v", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			if errors.Is(err, sql.ErrTxDone) {
				return
			}

			s.logger.Error(fmt.Sprintf("rollback: %s", err))
		}
	}()

	walletRepo := s.deps.NewWalletRepository(tx)

	// Создаем кошелек пользователю, если еще не создан.
	walletID, err := walletRepo.CreateIfNotExist(ctx, userID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create if not exist: %s", err))
		return 0, fmt.Errorf("create if not exist: %v", err)
	}

	// Зачисляем переданную сумму на кошелек пользователя.
	balance, err := walletRepo.Add(ctx, walletID, cash)
	if err != nil {
		s.logger.Error(fmt.Sprintf("add cash: %s", err))
		return 0, fmt.Errorf("add cash: %v", err)
	}

	// Генерируем payload.
	payload, err := transactions.EnrollmentPayload()
	if err != nil {
		s.logger.Error(fmt.Sprintf("generated payload: %s", err))
		return 0, fmt.Errorf("generated payload: %v", err)
	}

	txsRepo := s.deps.NewTransactionRepository(tx)

	// Добавляем транзакцию о проведеннной денежной операции.
	if err := txsRepo.AddTransaction(ctx, walletID, transactions.TypeAdd, payload, cash); err != nil {
		s.logger.Error(fmt.Sprintf("add transaction: %s", err))
		return 0, fmt.Errorf("add transaction: %v", err)
	}

	// Завершаем транзакцию.
	if err := tx.Commit(); err != nil {
		s.logger.Error(fmt.Sprintf("commit tx: %s", err))
		return 0, fmt.Errorf("commit tx: %v", err)
	}

	s.logger.Info(fmt.Sprintf("cash added for wallet: %d", walletID))

	return balance, nil
}
