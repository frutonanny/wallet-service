package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
)

const constraintName = "wallets_balance_check"

type Repository struct {
	db postgres.Database
}

func New(db postgres.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// ExistWallet - проверяет есть ли кошелек у пользователя, если есть возврает id кошелька.
// Если нет, то возвращаем ошибку ErrRepoWalletNotFound
func (r *Repository) ExistWallet(ctx context.Context, userID int64) (int64, error) {
	var walletID int64

	query := `select id from wallets where user_id = $1;`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repositories.ErrRepoWalletNotFound
		}
		return 0, fmt.Errorf("query row: %v", err)
	}

	return walletID, nil
}

// CreateIfNotExist - создает кошелек, если он ранее не был создан для переданного пользователя.
func (r *Repository) CreateIfNotExist(ctx context.Context, userID int64) (int64, error) {
	var walletID int64

	// Запрос взят https://stackoverflow.com/questions/40323799/return-rows-from-insert-with-on-conflict-without-needing-to-update
	query := `with ins as (
			insert into wallets (user_id) values ($1)
        		on conflict on constraint wallets_user_id_key do update
            		set user_id = null
            		where false
        		returning id)
			select id
			from ins
			union all
			select id
			from wallets
			where user_id = $1
			limit 1;`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	if err != nil {
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
func (r *Repository) Add(ctx context.Context, walletID int64, amount int64) (int64, error) {
	var balance int64

	query := `update wallets set balance = balance + $1 where id= $2 returning balance;`

	err := r.db.QueryRowContext(ctx, query, amount, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return balance, nil
}

// Reserve - резервирует переданную сумму денег.
func (r *Repository) Reserve(ctx context.Context, walletID, cash int64) (int64, error) {
	var balance int64

	query := `update wallets set reservation = reservation + $1, balance = balance - $1
				where id = $2 returning balance;`

	var pgErr *pgconn.PgError
	err := r.db.QueryRowContext(ctx, query, cash, walletID).Scan(&balance)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.ConstraintName == constraintName {
			return 0, repositories.ErrRepoNotEnoughCash
		}

		return 0, fmt.Errorf("query row: %w", err)
	}

	return balance, nil
}

// WriteOff - списывает переданную сумму денег
func (r *Repository) WriteOff(ctx context.Context, walletID, amount, delta int64) (int64, error) {
	var balance int64

	query := `update wallets set reservation = reservation - $1, balance = balance + $2
				where id = $3 returning balance;`

	err := r.db.QueryRowContext(ctx, query, amount, delta, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query row: %w", err)
	}

	return balance, nil
}

// Cancel - разрезервирует переданную сумму денег
func (r *Repository) Cancel(ctx context.Context, walletID, cash int64) (int64, error) {
	var balance int64

	query := `update wallets set reservation = reservation - $1, balance = balance + $1
				where id = $2 returning balance;`

	err := r.db.QueryRowContext(ctx, query, cash, walletID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("query row: %w", err)
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
