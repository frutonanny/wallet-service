//go:generate mockgen --source=service.go --destination=mock/service.go
package cancel

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

type WalletRepository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
	Cancel(ctx context.Context, walletID, cash int64) (int64, error)
}

type OrderRepository interface {
	GetOrder(ctx context.Context, externalID int64) (int64, string, int64, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) error
	AddOrderTransactions(ctx context.Context, orderID int64, nameType string) (int64, error)
}

type TransactionRepository interface {
	AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) (int64, error)
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewWalletRepository(db postgres.Database) WalletRepository
	NewOrderRepository(db postgres.Database) OrderRepository
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

// Cancel - разрезервирует переданную сумму средств у пользователя.
// - проверяем есть ли кошелек у пользователя, если нет, то отдаем ошибку ErrWalletNotFound.
// - проверяем есть ли заказ с переданным идентификатором внешнего заказа.
// 		1. Если заказа нет, то возвращаем ошибку ErrOrderNotFound.
// 		2. Заказ есть, то проверяем статус заказа. Должен быть reservation. Узнаем сумма резервирования.
// - списываем зарезервированную сумму с резерва пользователя. Добавляем эту сумму в баланс пользователя
// - обновляем информацию о заказе.
// - добавляем транзакцию об обновленном заказе;
// - добавляем транзакцию об отмене резервирования средств;
// - в ответ отдаем обновленный баланс пользователя в копейках.
func (s *Service) Cancel(ctx context.Context, userID, externalID int64) (int64, error) {
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

	// Проверяем, есть ли кошелек у пользователя.
	walletID, err := walletRepo.ExistWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoWalletNotFound) {
			return 0, servicesErrors.ErrWalletNotFound
		}

		s.logger.Error(fmt.Sprintf("wallet not exist: %s", err))
		return 0, fmt.Errorf("wallet not exist: %v", err)
	}

	orderRepo := s.deps.NewOrderRepository(tx)

	// Проверяем есть ли заказ с переданным идентификатором внешнего заказа.
	orderID, status, amount, err := orderRepo.GetOrder(ctx, externalID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoOrderNotFound) {
			return 0, servicesErrors.ErrOrderNotFound
		}

		s.logger.Error(fmt.Sprintf("order exist: %s", err))
		return 0, fmt.Errorf("order exist: %v", err)
	}

	// Проверяем соответствие статуса. Отменить резерв можно только в том случае, если заказ в резерве.
	if ok := orders.IsOrderReserved(status); !ok {
		s.logger.Error(fmt.Sprintf("order has wrong status %v", status))
		return 0, fmt.Errorf("order has wrong status %v", status)
	}

	//	Обновляем информацию в заказе.
	if err := orderRepo.UpdateOrderStatus(ctx, orderID, orders.StatusCancelled); err != nil {
		s.logger.Error(fmt.Sprintf("update order: %s", err))
		return 0, fmt.Errorf("update order: %v", err)
	}

	// Добавляем транзакцию об изменении статуса заказа.
	if _, err := orderRepo.AddOrderTransactions(ctx, orderID, orders.StatusCancelled); err != nil {
		s.logger.Error(fmt.Sprintf("add order transaction: %s", err))
		return 0, fmt.Errorf("add order transaction: %v", err)
	}

	// Разрезервируем переданную сумму. Эту сумму возвращаем в баланс.
	balance, err := walletRepo.Cancel(ctx, walletID, amount)
	if err != nil {
		s.logger.Error(fmt.Sprintf("cancel: %s", err))
		return 0, fmt.Errorf("cancel: %v", err)
	}

	// Генерируем payload.
	payload, err := transactions.CancelPayload(externalID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("generated payload: %s", err))
		return 0, fmt.Errorf("generated payload: %v", err)
	}

	txsRepo := s.deps.NewTransactionRepository(tx)

	// Добавляем транзакцию о разрезервированных средствах
	if _, err := txsRepo.AddTransaction(ctx, walletID, transactions.TypeCancel, payload, amount); err != nil {
		s.logger.Error(fmt.Sprintf("add transaction: %s", err))
		return 0, fmt.Errorf("add transaction: %v", err)
	}

	// Завершаем транзакцию.
	if err := tx.Commit(); err != nil {
		s.logger.Error(fmt.Sprintf("commit tx: %s", err))
		return 0, fmt.Errorf("commit tx: %v", err)
	}

	s.logger.Info(fmt.Sprintf("canceled cash reservation for wallet: %d", walletID))

	return balance, nil
}
