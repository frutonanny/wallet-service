package wallet

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	"github.com/frutonanny/wallet-service/internal/repositories"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
)

const (
	fileConfig = "../../../config/config.local.json"
)

var (
	config       = serviceConfig.Must(fileConfig)
	testUserID   = int64(10)
	testWalletID = int64(1)
	testCash     = int64(1_000)
	testDelta    = int64(0)
	data         = struct {
		Type string `json:"type"`
	}{
		Type: "enrollment",
	}
)

func TestCreateWallet(t *testing.T) {
	t.Run("create wallet successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)
	})
}

func TestExitWallet(t *testing.T) {
	t.Run("checked exist wallet successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		newWalletID, err := walletRepo.ExistWallet(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, walletID, newWalletID)
	})

	t.Run("checked exist wallet failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.ExistWallet(context.Background(), testUserID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, repositories.ErrRepoWalletNotFound)
	})
}

func TestCreateIfNotExist(t *testing.T) {
	t.Run("create wallet if not exist successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		walletID2, err := walletRepo.CreateIfNotExist(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, walletID, walletID2)
	})

	t.Run("create wallet if not exist successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.CreateIfNotExist(context.Background(), testUserID)
		assert.NoError(t, err)
	})
}

func TestAdd(t *testing.T) {
	t.Run("added amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		balance, err := walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, testCash, balance)
	})

	t.Run("added amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.Add(context.Background(), testWalletID, testCash)
		assert.Error(t, err)
	})
}

func TestReserve(t *testing.T) {
	t.Run("reserve amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		balance1, err := walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		balance2, err := walletRepo.Reserve(context.Background(), walletID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, balance1-testCash, balance2)
	})

	t.Run("reserve amount failed, not enough cash", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		_, err = walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		_, err = walletRepo.Reserve(context.Background(), walletID, 2*testCash)
		assert.Error(t, err)
		assert.ErrorIs(t, repositories.ErrRepoNotEnoughCash, err)
	})

	t.Run("reserve amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.Reserve(context.Background(), testWalletID, testCash)
		assert.Error(t, err)
	})
}

func TestWriteOff(t *testing.T) {
	t.Run("write-off amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		balance1, err := walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		balance2, err := walletRepo.Reserve(context.Background(), walletID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, balance1-testCash, balance2)

		balance3, err := walletRepo.WriteOff(context.Background(), walletID, testCash, testDelta)
		assert.NoError(t, err)
		assert.Equal(t, balance2, balance3)
	})

	t.Run("write-off amount failed, not enough reserved cash", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		_, err = walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		_, err = walletRepo.Reserve(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		_, err = walletRepo.WriteOff(context.Background(), walletID, 2*testCash, testDelta)
		assert.Error(t, err)
	})

	t.Run("write-off amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.WriteOff(context.Background(), testWalletID, testCash, testDelta)
		assert.Error(t, err)
	})
}

func TestCancel(t *testing.T) {
	t.Run("cancel amount successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		balance1, err := walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		balance2, err := walletRepo.Reserve(context.Background(), walletID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, balance1-testCash, balance2)

		balance3, err := walletRepo.Cancel(context.Background(), walletID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, balance1, balance3)
	})

	t.Run("cancel amount failed, not enough reserved cash", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		_, err = walletRepo.Add(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		_, err = walletRepo.Reserve(context.Background(), walletID, testCash)
		assert.NoError(t, err)

		_, err = walletRepo.Cancel(context.Background(), walletID, 2*testCash)
		assert.Error(t, err)
	})

	t.Run("cancel amount failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.Cancel(context.Background(), testWalletID, testCash)
		assert.Error(t, err)
	})
}

func TestGetBalance(t *testing.T) {
	t.Run("get balance successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		balance, err := walletRepo.GetBalance(context.Background(), walletID)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), balance)
	})

	t.Run("get balance failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		_, err := walletRepo.GetBalance(context.Background(), testWalletID)
		assert.Error(t, err)
	})
}

//query := []string{`
//					insert into transactions(wallet_id, "type", payload, amount)
//					values(10, typeIncoming, '{ "type": "enrollment" }', 100);`,
//}
