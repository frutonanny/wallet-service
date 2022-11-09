package wallet_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	"github.com/frutonanny/wallet-service/internal/postgres"
	"github.com/frutonanny/wallet-service/internal/repositories"
	repoWallet "github.com/frutonanny/wallet-service/internal/repositories/wallet"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
)

const (
	fileConfig   = "../../../config/config.local.json"
	testUserID   = int64(10)
	testWalletID = int64(1)
	testAmount   = int64(1_000)
	testDelta    = int64(0)
)

var config = serviceConfig.Must(fileConfig)

func TestRepository_CreateWallet(t *testing.T) {
	ctx := context.Background()
	t.Run("create wallet successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)
		assert.NotEmpty(t, walletID)

		// Проверяем, что кошелек создался.
		walletID2, err := walletRepo.ExistWallet(ctx, testUserID)
		require.NoError(t, err)
		assert.EqualValues(t, walletID, walletID2)
	})
}

func TestRepository_ExistWallet(t *testing.T) {
	ctx := context.Background()
	t.Run("checked exist wallet successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)
		assert.NotEmpty(t, walletID)

		// Проверяем, что кошелек создался.
		walletID2, err := walletRepo.ExistWallet(ctx, testUserID)
		require.NoError(t, err)
		assert.EqualValues(t, walletID, walletID2)
	})

	t.Run("checked exist wallet failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		_, err := walletRepo.ExistWallet(ctx, testUserID)
		require.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrRepoWalletNotFound)
	})
}

func TestRepository_CreateIfNotExist(t *testing.T) {
	ctx := context.Background()
	t.Run("wallet already created", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Создаем вновь кошелек пользователю, для которого создали ранее. Ожидаем, что запись не обновится
		// и получим тот же walletID.
		walletID2, err := walletRepo.CreateIfNotExist(ctx, testUserID)
		require.NoError(t, err)
		assert.EqualValues(t, walletID, walletID2)
	})

	t.Run("create wallet if not exist ", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек, так как ранее кошелек не был создан.
		walletID, err := walletRepo.CreateIfNotExist(ctx, testUserID)
		require.NoError(t, err)

		// Проверяем, что кошелек создался.
		walletID2, err := walletRepo.ExistWallet(ctx, testUserID)
		require.NoError(t, err)
		assert.EqualValues(t, walletID, walletID2)
	})
}

func TestRepository_Add(t *testing.T) {
	ctx := context.Background()
	t.Run("added amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount. В ответ получаем баланс = testAmount.
		balance, err := walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)
		assert.EqualValues(t, testAmount, balance)
	})

	t.Run("added amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := walletRepo.Add(ctx, testWalletID, testAmount)
		assert.Error(t, err)
	})
}

func TestRepository_Reserve(t *testing.T) {
	ctx := context.Background()
	t.Run("reserve amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount.
		balance1, err := walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Резервируем сумму testAmount. В ответ получаем измененный баланс balance1-testAmount
		balance2, err := walletRepo.Reserve(ctx, walletID, testAmount)
		require.NoError(t, err)
		assert.EqualValues(t, balance1-testAmount, balance2)

		// Проверяем зарезервированную сумму.
		reservation := getReservation(ctx, t, tx, walletID)
		assert.EqualValues(t, testAmount, reservation)
	})

	t.Run("reserve amount failed, not enough cash", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount.
		_, err = walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Резервируем сумму testAmount. В ответ получаем ошибку ErrRepoNotEnoughCash.
		_, err = walletRepo.Reserve(ctx, walletID, 2*testAmount)
		require.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrRepoNotEnoughCash)
	})

	t.Run("reserve amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := walletRepo.Reserve(ctx, testWalletID, testAmount)
		assert.Error(t, err)
	})
}

func TestRepository_WriteOff(t *testing.T) {
	ctx := context.Background()
	t.Run("write-off amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount.
		balance1, err := walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Резервируем сумму testAmount. В ответ получаем измененный баланс balance1-testAmount.
		balance2, err := walletRepo.Reserve(ctx, walletID, testAmount)
		require.NoError(t, err)
		assert.EqualValues(t, balance1-testAmount, balance2)

		// Проверяем зарезервированную сумму.
		reservation := getReservation(ctx, t, tx, walletID)
		assert.EqualValues(t, testAmount, reservation)

		// Списываем зарезервированную сумму. Баланс не должен измениться (разница между зарезервированным ранее
		// и списываемым testDelta=0)
		balance3, err := walletRepo.WriteOff(ctx, walletID, testAmount, testDelta)
		require.NoError(t, err)
		assert.EqualValues(t, balance2, balance3)

		// Проверяем зарезервированную сумму. Она должна быть = 0.
		reservation2 := getReservation(ctx, t, tx, walletID)
		assert.EqualValues(t, 0, reservation2)
	})

	t.Run("write-off amount failed, not enough reserved cash", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount.
		_, err = walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Резервируем сумму testAmount. В ответ получаем измененный баланс balance1-testAmount.
		_, err = walletRepo.Reserve(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Списываем бОльшую сумму, чем зарезервировано. Ожидаем ошибку.
		_, err = walletRepo.WriteOff(ctx, walletID, 2*testAmount, testDelta)
		assert.Error(t, err)
	})

	t.Run("write-off amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := walletRepo.WriteOff(ctx, testWalletID, testAmount, testDelta)
		assert.Error(t, err)
	})
}

func TestRepository_Cancel(t *testing.T) {
	ctx := context.Background()
	t.Run("cancel amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Добавляем на баланс сумму testAmount.
		balance1, err := walletRepo.Add(ctx, walletID, testAmount)
		require.NoError(t, err)

		// Резервируем сумму testAmount. В ответ получаем измененный баланс balance1-testAmount.
		balance2, err := walletRepo.Reserve(ctx, walletID, testAmount)
		require.NoError(t, err)
		assert.EqualValues(t, balance1-testAmount, balance2)

		// Проверяем зарезервированную сумму.
		reservation := getReservation(ctx, t, tx, walletID)
		assert.EqualValues(t, testAmount, reservation)

		// Разрезервируем сумму testAmount.
		balance3, err := walletRepo.Cancel(ctx, walletID, testAmount)
		require.NoError(t, err)
		assert.EqualValues(t, balance1, balance3)

		// Проверяем, что сумму разрезервирована. Она должна быть = 0.
		reservation2 := getReservation(ctx, t, tx, walletID)
		assert.EqualValues(t, 0, reservation2)
	})

	t.Run("cancel amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := walletRepo.Cancel(ctx, testWalletID, testAmount)
		assert.Error(t, err)
	})
}

func TestRepository_GetBalance(t *testing.T) {
	ctx := context.Background()
	t.Run("get balance successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Создаем кошелек.
		walletID, err := walletRepo.CreateWallet(ctx, testUserID)
		require.NoError(t, err)

		// Запрашиваем баланс, получаем дефолтное значение.
		balance, err := walletRepo.GetBalance(ctx, walletID)
		require.NoError(t, err)
		assert.EqualValues(t, 0, balance)
	})

	t.Run("get balance failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := repoWallet.New(tx)

		// Получаем ошибку, так как кошелька не существует.
		_, err := walletRepo.GetBalance(ctx, testWalletID)
		assert.Error(t, err)
	})
}

func getReservation(ctx context.Context, t *testing.T, db postgres.Database, walletID int64) int64 {
	t.Helper()

	var reservation int64

	query := `select reservation from wallets where id = $1;`

	err := db.QueryRowContext(ctx, query, walletID).Scan(&reservation)
	require.NoError(t, err)

	return reservation
}
