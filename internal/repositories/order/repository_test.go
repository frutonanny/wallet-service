package order_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	"github.com/frutonanny/wallet-service/internal/orders"
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoOrder "github.com/frutonanny/wallet-service/internal/repositories/order"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

const (
	fileConfig     = "../../../config/config.local.json"
	testUserID     = int64(1)
	testExternalID = int64(1)
	testServiceID  = int64(1)
	testAmount     = int64(500)
	testFailed     = int64(0)
)

var (
	config      = serviceConfig.Must(fileConfig)
	testStatusR = orders.StatusReserved
	testStatusW = orders.StatusWrittenOff
	testTypeTx  = transactions.TypeReserve
)

func TestRepository_CreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("create order successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)
		assert.NotEmpty(t, orderID)

		// Проверяем, что заказ создался c таким testExternalID создался.
		orderID2, status, amount, err := repo.GetOrder(ctx, testExternalID)
		require.NoError(t, err)
		assert.EqualValues(t, orderID, orderID2)
		assert.EqualValues(t, testStatusR, status)
		assert.EqualValues(t, testAmount, amount)
	})

	t.Run("create order failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := repo.CreateOrder(ctx, testFailed, testExternalID, testServiceID, testAmount, testStatusR)
		assert.Error(t, err)
	})
}

func TestRepository_GetOrderByServiceID(t *testing.T) {
	ctx := context.Background()

	t.Run("get order by serviceID successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)

		// Получаем информацию о ранее созданном заказе.
		orderID2, status, amount, err := repo.GetOrderByServiceID(ctx, testExternalID, testServiceID)
		require.NoError(t, err)
		assert.EqualValues(t, orderID, orderID2)
		assert.EqualValues(t, testStatusR, status)
		assert.EqualValues(t, testAmount, amount)
	})

	t.Run("get order by serviceID failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Получаем информацию о несуществующем заказе.
		_, _, _, err := repo.GetOrderByServiceID(ctx, testExternalID, testServiceID)
		assert.Error(t, err)
	})
}

func TestRepository_GetOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("get order successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)

		// Получаем информацию о ранее созданном заказе.
		orderID2, status, amount, err := repo.GetOrder(ctx, testExternalID)
		require.NoError(t, err)
		assert.EqualValues(t, orderID, orderID2)
		assert.EqualValues(t, testStatusR, status)
		assert.EqualValues(t, testAmount, amount)
	})

	t.Run("get order failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Получаем информацию о несуществующем заказе.
		_, _, _, err := repo.GetOrder(ctx, testExternalID)
		require.Error(t, err)
	})
}

func TestRepository_UpdateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("update order successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)

		// Обновляем статус заказа и стоимость уменьшаем в 2 раза.
		err = repo.UpdateOrder(ctx, orderID, testAmount/2, testStatusW)
		require.NoError(t, err)

		// Получаем информацию о ранее созданном заказе.
		orderID2, status, amount, err := repo.GetOrder(ctx, testExternalID)
		require.NoError(t, err)
		assert.EqualValues(t, orderID, orderID2)
		assert.EqualValues(t, testStatusW, status)
		assert.EqualValues(t, testAmount/2, amount)
	})

	t.Run("update order failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Обновляем информацию о несуществующем заказе.
		err := repo.UpdateOrder(ctx, testFailed, testAmount/2, testStatusW)
		assert.Error(t, err)
	})
}

func TestRepository_UpdateOrderStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("update order status successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)

		// Обновляем статус заказа.
		err = repo.UpdateOrderStatus(ctx, orderID, testStatusW)
		require.NoError(t, err)

		// Получаем информацию о ранее созданном заказе.
		orderID2, status, _, err := repo.GetOrder(ctx, testExternalID)
		require.NoError(t, err)
		assert.EqualValues(t, orderID, orderID2)
		assert.EqualValues(t, testStatusW, status)
	})

	t.Run("update order status failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Обновляем информацию о несуществующем заказе.
		err := repo.UpdateOrderStatus(ctx, testFailed, testStatusW)
		assert.Error(t, err)
	})
}

func TestRepository_AddOrderTransactions(t *testing.T) {
	ctx := context.Background()

	t.Run("add order transactions successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Создаем кошелек.
		walletID := createWallet(ctx, t, tx, testUserID)

		// Создаем заказ.
		orderID, err := repo.CreateOrder(ctx, walletID, testExternalID, testServiceID, testAmount, testStatusR)
		require.NoError(t, err)

		// Добавляем информацию о транзакции
		txID, err := repo.AddOrderTransactions(ctx, orderID, testTypeTx)
		require.NoError(t, err)

		txOrder := getOrderTx(ctx, t, tx, txID)
		assert.EqualValues(t, txID, txOrder.ID)
		assert.EqualValues(t, orderID, txOrder.OrderID)
		assert.EqualValues(t, testTypeTx, txOrder.Type)
	})

	t.Run("add order transactions failed, unknown walletID", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoOrder.New(tx)

		// Добавляем информацию о транзакции с несуществующим заказом.
		err := repo.UpdateOrderStatus(ctx, testFailed, testTypeTx)
		assert.Error(t, err)
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

type orderTx struct {
	ID      int64
	OrderID int64
	Type    string
}

// getOrderTx отдает информацию о транзакции.
func getOrderTx(ctx context.Context, t *testing.T, db postgres.Database, txID int64) orderTx {
	t.Helper()

	tx := orderTx{}

	query := `select id, order_id, "type" from order_transactions where id = $1;`

	err := db.QueryRowContext(ctx, query, txID).Scan(&tx.ID, &tx.OrderID, &tx.Type)
	require.NoError(t, err)

	return tx
}
