package order

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

// GetOrderByServiceID проверяет есть ли заказ с переданным идентификатором внешнего заказа и возращает информацию
// о заказе.
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
// идентификатор, статус заказа и его стоимость.
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

	res, err := r.db.ExecContext(ctx, query, status, amount, orderID) // как проверить обновилось поле или нет
	if err != nil {
		return fmt.Errorf("query row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affect: %w", err)
	}

	if rowsAffected == 0 {
		return repositories.ErrRepoOrderNotFound
	}

	return nil
}

// UpdateOrderStatus меняет статус заказа.
func (r *Repository) UpdateOrderStatus(ctx context.Context, orderID int64, status string) error {
	query := `update orders set status = $1 where id = $2;`

	res, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return fmt.Errorf("query row: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affect: %w", err)
	}

	if rowsAffected == 0 {
		return repositories.ErrRepoOrderNotFound
	}

	return nil
}

// AddOrderTransactions - добавляет информацию о действиях с заказами.
func (r *Repository) AddOrderTransactions(ctx context.Context, orderID int64, nameType string) (int64, error) {
	var id int64

	query := `insert into order_transactions(order_id, "type") values($1, $2) returning id;`

	err := r.db.QueryRowContext(ctx, query, orderID, nameType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("query row context: %v", err)
	}

	return id, nil
}
