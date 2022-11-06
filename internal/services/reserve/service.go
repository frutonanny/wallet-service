//go:generate mockgen --source=service.go --destination=mock/service.go
package reserve

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/frutonanny/wallet-service/internal/orders"
	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

type logger interface {
	Info(msg string)
	Error(msg string)
}

type walletRepository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
	Reserve(ctx context.Context, walletID, cash int64) (int64, error)
	CreateOrder(ctx context.Context, walletID, externalID, serviceID, amount int64, status string) (int64, error)
	AddOrderTransactions(ctx context.Context, orderID int64, nameType string) error
}

type transactionRepository interface {
	AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewWalletRepository(db postgres.Database) walletRepository
	NewTransactionRepository(db postgres.Database) transactionRepository
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

// Reserve - резервирует переданную сумму средст у пользователя для оплаты заказа.
// - проверяем есть ли кошелек у пользователя, если нет, то отдаем ошибку ErrWalletNotFound.
// - проверяем достаточно ли средств у пользователя, если нет, то возвращаем ошибку ErrNotEnoughCash.
// - списываем переданную сумму с баланса пользователя и добавляем эту сумму в резерв.
// - создаем заказ со статусом "reservation".
// - добавляем транзакцию о созданном заказе;
// - добавляем транзакцию о зарезервированных средствах;
// - в ответ отдаем обновленный баланс пользователя в копейках, без учета зерезервированных денег.
func (s *Service) Reserve(ctx context.Context, userID, serviceID, externalID, price int64) (int64, error) {
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

	// Проверяем, есть ди кошелек у пользователя.
	walletID, err := walletRepo.ExistWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoWalletNotFound) {
			return 0, servicesErrors.ErrWalletNotFound
		}

		s.logger.Error(fmt.Sprintf("wallet not exist: %s", err))
		return 0, fmt.Errorf("wallet not exist: %v", err)
	}

	// Списываем переданную сумму с баланса пользователя и добавляем эту сумму в резерв.
	// Одновременно проверяем достаточно ли средств у пользователя.
	balance, err := walletRepo.Reserve(ctx, walletID, price)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoNotEnoughCash) {
			return 0, servicesErrors.ErrNotEnoughCash
		}

		s.logger.Error(fmt.Sprintf("reserve: %s", err))
		return 0, fmt.Errorf("reserve: %v", err)
	}

	// Создаем заказ со статусом "reservation".
	orderID, err := walletRepo.CreateOrder(ctx, walletID, externalID, serviceID, price, orders.StatusReserved)
	if err != nil {
		s.logger.Error(fmt.Sprintf("create order: %s", err))
		return 0, fmt.Errorf("create order: %v", err)
	}

	// Добавляем транзакцию о созданном заказе
	if err := walletRepo.AddOrderTransactions(ctx, orderID, orders.StatusReserved); err != nil {
		s.logger.Error(fmt.Sprintf("add order transaction: %s", err))
		return 0, fmt.Errorf("add order transaction: %v", err)
	}

	// Генерируем payload.
	payload, err := transactions.ReservationPayload(externalID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("generated payload: %s", err))
		return 0, fmt.Errorf("generated payload: %v", err)
	}

	txsRepo := s.deps.NewTransactionRepository(tx)

	// Добавляем транзакцию о зарезервированных средствах
	if err := txsRepo.AddTransaction(ctx, walletID, transactions.TypeReserve, payload, price); err != nil {
		s.logger.Error(fmt.Sprintf("add transaction: %s", err))
		return 0, fmt.Errorf("add transaction: %v", err)
	}

	// Завершаем транзакцию.
	if err := tx.Commit(); err != nil {
		s.logger.Error(fmt.Sprintf("commit tx: %s", err))
		return 0, fmt.Errorf("commit tx: %v", err)
	}

	s.logger.Info(fmt.Sprintf("cash reserved for wallet: %d", walletID))

	return balance, nil
}
