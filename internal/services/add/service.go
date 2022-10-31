//go:generate mockgen --source=service.go --destination=mock/service.go
package add

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
)

const actionName = "incoming_transfer"

var Data = struct {
	Type string `json:"type"`
}{
	Type: "enrollment",
}

type logger interface {
	Info(msg string)
	Error(msg string)
}

type Repository interface {
	ExistWallet(ctx context.Context, userID int64) (int64, error)
	CreateWallet(ctx context.Context, userID int64) (int64, error)
	Add(ctx context.Context, walletID int64, cash int64) (int64, error)
	AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error
}

// RepoBuilder умеет налету создавать репозиторий поверх *sql.DB, *sql.Tx.
// Нужен для написания юнит-тестов без подключения к базе.
type RepoBuilder interface {
	NewRepository(db postgres.Database) Repository
}

type Service struct {
	db      *sql.DB
	logger  logger
	builder RepoBuilder
}

func New(logger logger, db *sql.DB) *Service {
	return &Service{
		logger:  logger,
		db:      db,
		builder: &builderImpl{},
	}
}

func (s *Service) WithBuilder(builder RepoBuilder) *Service {
	s.builder = builder
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

	repo := s.builder.NewRepository(tx)

	var walletID int64
	// Проверяем есть ли кошелек у пользователя.
	walletID, err = repo.ExistWallet(ctx, userID)
	if err != nil {
		if !errors.Is(err, repositories.ErrRepoWalletNotFound) {
			s.logger.Error(fmt.Sprintf("exist wallet: %s", err))
			return 0, fmt.Errorf("exist wallet: %v", err)
		}
		// Создаем кошелек
		walletID, err = repo.CreateWallet(ctx, userID)
		if err != nil {
			s.logger.Error(fmt.Sprintf("create wallet: %s", err))
			return 0, fmt.Errorf("create wallet: %v", err)
		}
	}

	balance, err := repo.Add(ctx, walletID, cash)
	if err != nil {
		s.logger.Error(fmt.Sprintf("add cash: %s", err))
		return 0, fmt.Errorf("add cash: %v", err)
	}

	// TODO
	payload, err := json.Marshal(Data)
	if err != nil {
		s.logger.Error(fmt.Sprintf("marshal payloud: %s", err))
		return 0, fmt.Errorf("marshal payloud: %v", err)
	}

	// Добавляем транзакцию о проведеннной денежной операции.
	if err := repo.AddTransaction(ctx, walletID, actionName, payload, cash); err != nil {
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
