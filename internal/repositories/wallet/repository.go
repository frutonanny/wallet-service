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
// Если нет, тозвращаем ошибку ErrRepoWalletNotFound
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
func (r *Repository) Add(ctx context.Context, walletID int64, cash int64) (int64, error) {
	var balance int64

	query := `update wallets set balance = balance + $1 where id= $2 returning balance;`

	err := r.db.QueryRowContext(ctx, query, cash, walletID).Scan(&balance)
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

// CreateOrder создает заказ.
func (r *Repository) CreateOrder(
	ctx context.Context,
	walletID,
	externalID,
	serviceID,
	amount int64,
	status string,
) (int64, error) {
	var orderID int64

	query := `insert into orders(wallet_id, external_id, service_id, status, amount) 
				values($1, $2, $3, $4, $5) returning id;`

	err := r.db.QueryRowContext(ctx, query, walletID, externalID, serviceID, status, amount).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return orderID, nil
}

// GetOrderByServiceID проверяет есть ли заказ с переданным идентификатором внешнего заказа и возращает информацию о заказе.
// Если заказа нет, то возвращаем ошибку ErrRepoOrderNotFound.
func (r *Repository) GetOrderByServiceID(
	ctx context.Context,
	externalID, serviceID int64,
) (int64, string, int64, error) {
	var orderID, amount int64
	var status string

	query := `select id, status, amount from orders where external_id = $1 and service_id = $2;`

	err := r.db.QueryRowContext(ctx, query, externalID, serviceID).Scan(&orderID, &status, &amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", 0, repositories.ErrRepoOrderNotFound
		}
		return 0, "", 0, fmt.Errorf("query row: %v", err)
	}

	return orderID, status, amount, nil
}

// GetOrder проверяет есть ли заказ с переданным идентификатором внешнего заказа и возращает
// идентификатор и статус заказа.
// Если заказа нет, то возвращаем ошибку ErrRepoOrderNotFound.
func (r *Repository) GetOrder(ctx context.Context, externalID int64) (int64, string, int64, error) {
	var orderID, amount int64
	var status string

	query := `select id, status, amount from orders where external_id = $1;`

	err := r.db.QueryRowContext(ctx, query, externalID).Scan(&orderID, &status, &amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, "", 0, repositories.ErrRepoOrderNotFound
		}
		return 0, "", 0, fmt.Errorf("query row: %v", err)
	}

	return orderID, status, amount, nil
}

// UpdateOrder меняет описание заказа (статус и стоимость).
func (r *Repository) UpdateOrder(ctx context.Context, orderID, amount int64, status string) error {
	query := `update orders set status = $1, amount = $2 where id = $3;`

	_, err := r.db.ExecContext(ctx, query, status, amount, orderID)
	if err != nil {
		return fmt.Errorf("query row: %v", err)
	}

	return nil
}

// UpdateOrderStatus меняет статус заказа.
func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	query := `update orders set status = $1 where id = $3;`

	err := r.db.QueryRowContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("query row: %v", err)
	}

	return nil
}

// AddOrderTransactions - добавляет информацию о действиях с заказами.
func (r *Repository) AddOrderTransactions(ctx context.Context, orderID int64, nameType string) error {
	query := `insert into order_transactions(order_id, "type") values($1, $2);`

	_, err := r.db.ExecContext(ctx, query, orderID, nameType)
	if err != nil {
		return fmt.Errorf("exec query: %v", err)
	}

	return nil
}
