package transaction

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

const (
	fileConfig  = "../../../config/config.local.json"
	testUserID  = int64(1)
	testAmount  = int64(500)
	testOrderID = int64(5)
	testLimit   = int64(5)
	testOffset  = int64(1)
)

var (
	config = serviceConfig.Must(fileConfig)
)

func TestGetTransactions(t *testing.T) {
	// Генерируем payload Enrollment.
	payloadAdd, err := transactions.EnrollmentPayload()
	assert.NoError(t, err)

	// Генерируем payload Reservation.
	payloadReserv, err := transactions.ReservationPayload(testOrderID)
	assert.NoError(t, err)

	t.Run("get transactions successfully", func(t *testing.T) {
		ctx := context.Background()

		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := New(tx)

		// Создаем кошелек.
		walletID, err := repo.createWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем транзакцио о пополнении.
		err = repo.AddTransaction(ctx, walletID, transactions.TypeAdd, payloadAdd, testAmount)
		require.NoError(t, err)

		// Добавляем транзакцио о резервировании
		err = repo.AddTransaction(ctx, walletID, transactions.TypeReserve, payloadReserv, testAmount)
		require.NoError(t, err)

		result, err := repo.GetTransactions(ctx, walletID, testLimit, testOffset, Date, Asc)
		require.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, transactions.TypeReserve, result[0].Type)
	})

	t.Run("get transactions successfully", func(t *testing.T) {
		ctx := context.Background()

		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := New(tx)

		// Создаем кошелек.
		walletID, err := repo.createWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем транзакцио о пополнении.
		err = repo.AddTransaction(ctx, walletID, transactions.TypeAdd, payloadAdd, 3*testAmount)
		require.NoError(t, err)

		// Добавляем транзакцио о пополнении.
		err = repo.AddTransaction(ctx, walletID, transactions.TypeAdd, payloadAdd, 2*testAmount)
		require.NoError(t, err)

		// Добавляем транзакцио о резервировании
		err = repo.AddTransaction(ctx, walletID, transactions.TypeReserve, payloadReserv, testAmount)
		require.NoError(t, err)

		result, err := repo.GetTransactions(ctx, walletID, testLimit, testOffset, Amount, Desc)
		require.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, transactions.TypeReserve, result[1].Type)
	})
}

// createWallet создает кошелек. Метод нужен для тестирования
func (r *Repository) createWallet(ctx context.Context, userID int64) (int64, error) {
	var walletID int64

	query := `insert into wallets(user_id) values($1) returning id;`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	if err != nil {
		return 0, fmt.Errorf("query row: %v", err)
	}

	return walletID, nil
}

// TODO тест на byTime
