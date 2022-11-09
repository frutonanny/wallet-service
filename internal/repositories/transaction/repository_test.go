package transaction_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoTxs "github.com/frutonanny/wallet-service/internal/repositories/transaction"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

const (
	fileConfig   = "../../../config/config.local.json"
	testWalletID = int64(52)
	testUserID   = int64(1)
	testAmount   = int64(500)
	testLimit    = int64(4)
	testOffset   = int64(1)
	testFailed   = int64(0)
)

var (
	config = serviceConfig.Must(fileConfig)
)

func TestRepository_AddTransaction(t *testing.T) {
	ctx := context.Background()

	// Генерируем payload Enrollment.
	payloadAdd, err := transactions.EnrollmentPayload()
	assert.NoError(t, err)

	t.Run("add transactions successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoTxs.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Добавляем транзакцию о пополнении баланса.
		txID, err := repo.AddTransaction(ctx, walletID, transactions.TypeAdd, payloadAdd, testAmount)
		require.NoError(t, err)

		// Получаем данные транзакции для проверки
		txAdd := getTx(ctx, t, tx, txID)

		// Проверяем, что транзакция добавлена. Проверяем значения.
		assert.EqualValues(t, txID, txAdd.ID)
		assert.EqualValues(t, walletID, txAdd.WalletID)
		assert.EqualValues(t, transactions.TypeAdd, txAdd.Type)
		assert.EqualValues(t, testAmount, txAdd.Amount)
	})

	t.Run("add transactions failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := repo.AddTransaction(ctx, testFailed, transactions.TypeAdd, payloadAdd, 3*testAmount)
		assert.Error(t, err)
	})
}

func TestRepository_GetTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("get transactions by date descending successfully", func(t *testing.T) {

		query := []string{`insert into wallets(id, user_id, balance) 
							values(52, 7, 5000);`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 1000, '2022-11-01 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 2000, '2022-11-02 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 3000, '2022-11-03 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 4000, '2022-11-04 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 5000, '2022-11-05 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 6000, '2022-11-06 12:00');`,
		}

		// Заполняем предварительно данными, для корректной выборки.
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN, query)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получим список транзакций, отсортированный по убыванию даты с ограничениями: смещение = 1, лимит = 4.
		result, err := repo.GetTransactions(ctx, testWalletID, testLimit, testOffset, repoTxs.Date, repoTxs.Desc)
		require.NoError(t, err)
		assert.Len(t, result, 4)

		// транзакция от 2022-11-06 12:00 в выборку не попала из-за смещения
		assert.EqualValues(t, 5000, result[0].Amount) // транзакция от 2022-11-05 12:00
		assert.EqualValues(t, 2000, result[3].Amount) // транзакция от 2022-11-02 12:00
	})

	t.Run("get transactions in ascending amount successfully", func(t *testing.T) {

		query := []string{`insert into wallets(id, user_id, balance) 
							values(52, 7, 5000);`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 1000, '2022-11-01 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 2000, '2022-11-02 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 3000, '2022-11-03 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 4000, '2022-11-04 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 5000, '2022-11-05 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 6000, '2022-11-06 12:00');`,
		}

		// Заполняем предварительно данными, для корректной выборки.
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN, query)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получим список транзакций, отсортированный по убыванию даты с ограничениями: смещение = 1, лимит = 4.
		result, err := repo.GetTransactions(ctx, testWalletID, testLimit, testOffset, repoTxs.Amount, repoTxs.Asc)
		require.NoError(t, err)
		assert.Len(t, result, 4)

		// транзакция от 2022-11-01 12:00 в выборку не попала из-за смещения
		assert.EqualValues(t, 2000, result[0].Amount) // транзакция от 2022-11-02 12:00
		assert.EqualValues(t, 5000, result[3].Amount) // транзакция от 2022-11-05 12:00
	})

	t.Run("get transactions failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получаем пустой результат, так как кошелька не существует. TODO
		result, err := repo.GetTransactions(ctx, testFailed, testLimit, testOffset, repoTxs.Amount, repoTxs.Asc)
		require.NoError(t, err)
		require.Empty(t, result)
	})
}

func TestRepository_GetTransactionsByTime(t *testing.T) {
	ctx := context.Background()

	start, err := time.Parse(time.RFC3339, "2022-11-02T12:00:00Z")
	require.NoError(t, err)

	end, err := time.Parse(time.RFC3339, "2022-11-06T12:00:00Z")
	require.NoError(t, err)

	t.Run("get transactions by time successfully", func(t *testing.T) {

		query := []string{`insert into wallets(id, user_id, balance) 
							values(52, 7, 5000);`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 1000, '2022-11-01 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 2000, '2022-11-02 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 3000, '2022-11-03 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 4000, '2022-11-04 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'incoming_transfer', '{ "type": "enrollment" }', 5000, '2022-11-05 12:00');`,
			`insert into transactions(wallet_id, "type", payload, amount, created_at)
					values(52, 'reservation', '{ "order_id": 10 }', 6000, '2022-11-06 12:00');`,
		}

		// Заполняем предварительно данными, для корректной выборки.
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN, query)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получим список транзакций, отсортированный по убыванию даты, ограниченный датами start <= created_at <= end.
		result, err := repo.GetTransactionsByTime(ctx, testWalletID, start, end)
		require.NoError(t, err)
		assert.Len(t, result, 5)
		assert.EqualValues(t, 6000, result[0].Amount) // транзакция от 2022-11-06 12:00
		assert.EqualValues(t, 2000, result[4].Amount) // транзакция от 2022-11-02 12:00
	})

	t.Run("get transactions by time failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoTxs.New(tx)

		// Получаем пустой результат, так как кошелька не существует. TODO
		result, err := repo.GetTransactionsByTime(ctx, testFailed, start, end)
		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

// createWallet создает кошелек.
func createWallet(ctx context.Context, t *testing.T, db postgres.Database, userID int64) int64 {
	t.Helper()

	var walletID int64

	query := `insert into wallets(user_id) values($1) returning id;`

	err := db.QueryRowContext(ctx, query, userID).Scan(&walletID)
	require.NoError(t, err)

	return walletID
}

type transaction struct {
	ID       int64
	WalletID int64
	Type     string
	Amount   int64
}

// getOrderTx отдает информацию о транзакции.
func getTx(ctx context.Context, t *testing.T, db postgres.Database, txID int64) transaction {
	t.Helper()

	tx := transaction{}

	query := `select id, wallet_id, "type", amount from transactions where id = $1;`

	err := db.QueryRowContext(ctx, query, txID).Scan(&tx.ID, &tx.WalletID, &tx.Type, &tx.Amount)
	require.NoError(t, err)

	return tx
}
