package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
)

type Repository struct {
	db postgres.Database
}

func New(db postgres.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// ExistWallet - проверяет есть ли кошелек у пользователя, если есть возврает id кошелька.
func (r *Repository) ExistWallet(ctx context.Context, userID int64) (int64, error) {
	query := `select id from wallets where user_id = $1;`

	var walletID int64

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repositories.ErrRepoWalletNotFound
		}
		return 0, fmt.Errorf("query row: %v", err)
	}

	return walletID, nil
}

// CreateWallet - создает кошелек для переданного пользователя.
func (r *Repository) CreateWallet(ctx context.Context, userID int64) (int64, error) {
	var walletID int64

	query := `insert into wallets(user_id) values($1) returning id;`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return walletID, nil
}

// Add - зачисляет переданную сумму на кошелек пользователя и возвращает текущий баланс.
func (r *Repository) Add(ctx context.Context, walletID int64, cash int64) (int64, error) {
	query := `update wallets set balance = balance + $1 where id= $2 returning balance;`
	var balance int64

	err := r.db.QueryRowContext(ctx, query, cash, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return balance, nil
}

// GetBalance - отдает текущий баланс пользователя.
func (r *Repository) GetBalance(ctx context.Context, walletID int64) (int64, error) {
	var balance int64

	query := `select balance from wallets where id=$1;`

	err := r.db.QueryRowContext(ctx, query, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return balance, nil
}

// AddTransaction - добавляет информацию о проведеннной денежной операции.
func (r *Repository) AddTransaction(ctx context.Context, walletID int64, action string, payload []byte, amount int64) error {
	query := `insert into transactions(wallet_id, "type", payload, amount) values($1, $2, $3, $4);`

	_, err := r.db.ExecContext(ctx, query, walletID, action, payload, amount)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}
