package transaction

import (
	"context"
	"fmt"
	"time"

	"github.com/frutonanny/wallet-service/internal/postgres"
)

type Repository struct {
	db postgres.Database
}

func New(db postgres.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// AddTransaction - добавляет информацию о проведенной денежной операции.
func (r *Repository) AddTransaction(
	ctx context.Context,
	walletID int64,
	action string,
	payload []byte,
	amount int64,
) (int64, error) {
	var id int64

	query := `insert into transactions(wallet_id, "type", payload, amount) values($1, $2, $3, $4) returning id;`

	err := r.db.QueryRowContext(ctx, query, walletID, action, payload, amount).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("exec query: %v", err)
	}

	return id, nil
}

// GetTransactions отдает список транзакций пользователя, отсортированный по переданному параметру.
//
// Запрос с sortBy == "amount" потенциально тяжелый. Добавление индекса на колонку amount не имеет большого смысла
// из-за невысокой селективности значений. БД придется выбрать все транзакции по кошельку и отсортировать их.
// Как правило, банки предоставляют сортировку только по дате транзакций.
func (r *Repository) GetTransactions(
	ctx context.Context,
	walletID, limit, offset int64,
	sortBy SortBy,
	direction Direction,
) ([]Transaction, error) {
	query := `select "type", payload, amount, created_at
		from transactions
		where wallet_id = $1 order by ` + string(sortBy) + ` ` + string(direction) + ` limit $2 offset $3;`

	rows, err := r.db.QueryContext(ctx, query, walletID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var result []Transaction

	for rows.Next() {
		tr := Transaction{}

		err = rows.Scan(
			&tr.Type, &tr.Payload, &tr.Amount, &tr.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, tr)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}

// GetTransactionsByTime - выводит список транзакций для пользователя в переданном промежутке времени, отсортированный
// по переданному направлению.
func (r *Repository) GetTransactionsByTime(
	ctx context.Context,
	walletID int64,
	timeStart, timeEnd time.Time,
) ([]Transaction, error) {
	query := `select "type", payload, amount, created_at
		from transactions
		where wallet_id = $1 and (created_at >= $2 and created_at <= $3)
		order by created_at  desc;`

	rows, err := r.db.QueryContext(ctx, query, walletID, timeStart, timeEnd)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		_ = rows.Close()
	}()

	var result []Transaction

	for rows.Next() {
		tx := Transaction{}

		if err := rows.Scan(&tx.Type, &tx.Payload, &tx.Amount, &tx.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan: %w", err)
		}

		result = append(result, tx)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows: %w", err)
	}

	return result, nil
}
