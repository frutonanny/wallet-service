package wallet

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	"github.com/frutonanny/wallet-service/internal/repositories"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
)

const (
	fileConfig   = "../../../config/config.json"
	typeIncoming = "incoming_transfer"
)

var (
	config       = serviceConfig.Must(fileConfig)
	testUserID   = int64(1)
	testWalletID = int64(1)
	testCash     = int64(1_000)
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

func TestAddTransaction(t *testing.T) {
	t.Run("add transaction successfully", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		walletID, err := walletRepo.CreateWallet(context.Background(), testUserID)
		assert.NoError(t, err)

		payload, err := json.Marshal(data)
		assert.NoError(t, err)

		err = walletRepo.AddTransaction(context.Background(), walletID, typeIncoming, payload, testCash)
		assert.NoError(t, err)
	})

	t.Run("add transaction failed", func(t *testing.T) {
		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		walletRepo := New(tx)

		payload, err := json.Marshal(data)
		assert.NoError(t, err)

		err = walletRepo.AddTransaction(context.Background(), testWalletID, typeIncoming, payload, testCash)
		assert.Error(t, err)
	})
}

//query := []string{`
//					insert into transactions(wallet_id, "type", payload, amount)
//					values(10, typeIncoming, '{ "type": "enrollment" }', 100);`,
//}
