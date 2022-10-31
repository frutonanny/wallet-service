package testing_boilerplate

import (
	"context"
	"database/sql"
	"testing"

	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/stretchr/testify/require"
)

// InitDB - отдает соединение с базой, предварительно заполняея ее переданными данными queries-запросами.
func InitDB(t *testing.T, dsn string, queries ...[]string) (*sql.Tx, func()) {
	db, cancel := connectDB(dsn)

	tx := factory(t, db, queries...)

	cancel2 := func() {
		_ = tx.Rollback()
		cancel()
	}

	return tx, cancel2
}

// connectDB - устанавливает соединение с базой.
func connectDB(dsn string) (*sql.DB, func()) {
	db := postgres.MustConnect(dsn)

	cancel := func() {
		_ = db.Close()
	}

	return db, cancel
}

//factory - в транзакции выполеняет переданные sql-запросы, а также проверяет безошибочное их выполнение.
func factory(t *testing.T, db *sql.DB, queries ...[]string) *sql.Tx {
	ctx := context.Background()

	tx := beginTx(t, db)

	for i := range queries {
		for _, query := range queries[i] {
			_, err := tx.ExecContext(ctx, query)
			if err != nil {
				err2 := tx.Rollback()

				require.NoError(t, err)
				require.NoError(t, err2)
			}
		}
	}

	return tx
}

// beginTx - запускает транзакцию.
func beginTx(t *testing.T, db *sql.DB) *sql.Tx {
	tx, err := db.BeginTx(context.Background(), nil)
	require.NoError(t, err)
	return tx
}
