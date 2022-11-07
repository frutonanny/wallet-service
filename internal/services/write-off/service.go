package write_off

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	WriteOff(ctx context.Context, walletID, amount, delta int64) (int64, error)
}
type OrderRepository interface {
	GetOrderByServiceID(ctx context.Context, externalID, serviceID int64) (int64, string, int64, error)
	UpdateOrder(ctx context.Context, orderID, amount int64, status string) error
	AddOrderTransactions(ctx context.Context, orderID int64, nameType string) (int64, error)
}

type TransactionRepository interface {
	AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error
}

type ReportRepository interface {
	AddRecord(ctx context.Context, serviceID, amount int64, now time.Time) error
}

// dependencies умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type dependencies interface {
	NewWalletRepository(db postgres.Database) WalletRepository
	NewOrderRepository(db postgres.Database) OrderRepository
	NewTransactionRepository(db postgres.Database) TransactionRepository
	NewReportRepository(db postgres.Database) ReportRepository
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

// WriteOff - списывает переданную сумму средст у пользователя для оплаты заказа.
// - проверяем есть ли кошелек у пользователя, если нет, то отдаем ошибку ErrWalletNotFound.
// - проверяем есть ли заказ с переданным идентификатором внешнего заказа.
// 1. Если заказа нет, то возвращаем ошибку ErrOrderNotFound.
// 2. Заказ есть, то проверяем сумму списания и статус заказа, если price <= ранее переданной зарезервированной суммы,
// то списывваем новую переданную сумму с резерва пользователя, а разницу добавляем в баланс.
// - обновляем информацию о заказе.
// - добавляем транзакцию об обновленном заказе;
// - добавляем транзакцию о списанных средствах;
// - Записываем в отчет списание.
// - в ответ отдаем обновленный баланс пользователя в копейках.
func (s *Service) WriteOff(ctx context.Context, userID, serviceID, externalID, price int64) (int64, error) {
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

	orderRepo := s.deps.NewOrderRepository(tx)

	// Проверяем есть ли заказ с переданным идентификатором внешнего заказа.
	orderID, status, amount, err := orderRepo.GetOrderByServiceID(ctx, externalID, serviceID)
	if err != nil {
		if errors.Is(err, repositories.ErrRepoOrderNotFound) {
			return 0, servicesErrors.ErrOrderNotFound
		}

		s.logger.Error(fmt.Sprintf("order exist: %s", err))
		return 0, fmt.Errorf("order exist: %v", err)
	}

	// Проверяем соответствие статуса и стоимости.
	if ok := orders.IsOrderReserved(status); !ok {
		s.logger.Error(fmt.Sprintf("order has wrong status %v", status))
		return 0, fmt.Errorf("order has wrong status %v", status)
	}

	// Не даем списать бОльшую сумму, чем сейчас находится в резерве по этому заказу.
	if price > amount {
		s.logger.Error(fmt.Sprintf("the new value is greater than the reserved amount"))
		return 0, fmt.Errorf("the new value is greater than the reserved amount")
	}

	// Обновляем информацию в заказе.
	if err := orderRepo.UpdateOrder(ctx, orderID, price, orders.StatusWrittenOff); err != nil {
		s.logger.Error(fmt.Sprintf("update order: %s", err))
		return 0, fmt.Errorf("update order: %v", err)
	}

	// Добавляем транзакцию об изменении статуса заказа.
	if _, err := orderRepo.AddOrderTransactions(ctx, orderID, orders.StatusWrittenOff); err != nil {
		s.logger.Error(fmt.Sprintf("add order transaction: %s", err))
		return 0, fmt.Errorf("add order transaction: %v", err)
	}

	// Списываем переданную сумму с резерва. Если сумма списанная меньше зарезервированной, то возвращаем разницу в
	// баланс.
	balance, err := walletRepo.WriteOff(ctx, walletID, amount, amount-price)
	if err != nil {
		s.logger.Error(fmt.Sprintf("write-off: %s", err))
		return 0, fmt.Errorf("write-off: %v", err)
	}

	// Генерируем payload.
	payload, err := transactions.WriteOffPayload(externalID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("write-off payload: %s", err))
		return 0, fmt.Errorf("write-off payload: %v", err)
	}

	txsRepo := s.deps.NewTransactionRepository(tx)

	// Добавляем транзакцию о зарезервированных средствах.
	if err := txsRepo.AddTransaction(ctx, walletID, transactions.TypeWriteOff, payload, price); err != nil {
		s.logger.Error(fmt.Sprintf("add transaction: %s", err))
		return 0, fmt.Errorf("add transaction: %v", err)
	}

	reportRepo := s.deps.NewReportRepository(tx)

	if err := reportRepo.AddRecord(ctx, serviceID, amount, time.Now()); err != nil {
		s.logger.Error(fmt.Sprintf("add record: %s", err))
		return 0, fmt.Errorf("add record: %v", err)
	}

	// Завершаем транзакцию.
	if err := tx.Commit(); err != nil {
		s.logger.Error(fmt.Sprintf("commit tx: %s", err))
		return 0, fmt.Errorf("commit tx: %v", err)
	}

	s.logger.Info(fmt.Sprintf("cash written-off for wallet: %d", walletID))

	return balance, nil
}
