package get_transactions_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/repositories"
	"github.com/frutonanny/wallet-service/internal/repositories/transaction"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions"
	mock_get_txs "github.com/frutonanny/wallet-service/internal/services/get_transactions/mock"
)

const (
	testUserID   = int64(1)
	testWalletID = int64(1)
	testLimit    = int64(1)
	testOffset   = int64(1)
)

var testError = errors.New("error")

func TestService_GetTransactions(t *testing.T) {
	var db *sql.DB

	t.Run("get transaction successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		repoTxs := mock_get_txs.NewMockTransactionRepository(ctrl)
		repoTxs.
			EXPECT().
			GetTransactions(ctx, testWalletID, testLimit, testOffset, transaction.Amount, transaction.Desc).
			Return(nil, nil)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(repoTxs)

		log := mock_get_txs.NewMocklogger(ctrl)

		server := get_transactions.New(log, db).WithDependencies(deps)
		_, err := server.GetTransactions(
			ctx,
			testUserID,
			testLimit,
			testOffset,
			get_transactions.Amount,
			get_transactions.Desc)
		assert.NoError(t, err)
	})

	t.Run("get transaction failed, ErrWalletNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, repositories.ErrRepoWalletNotFound)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		server := get_transactions.New(log, db).WithDependencies(deps)
		_, err := server.GetTransactions(
			ctx,
			testUserID,
			testLimit,
			testOffset,
			get_transactions.Amount,
			get_transactions.Desc)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("get transaction failed, exist wallet error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, testError)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		server := get_transactions.New(log, db).WithDependencies(deps)
		_, err := server.GetTransactions(
			ctx,
			testUserID,
			testLimit,
			testOffset,
			get_transactions.Amount,
			get_transactions.Desc)
		assert.Error(t, err)
	})

	t.Run("get transaction failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		repoTxs := mock_get_txs.NewMockTransactionRepository(ctrl)
		repoTxs.
			EXPECT().
			GetTransactions(ctx, testWalletID, testLimit, testOffset, transaction.Amount, transaction.Desc).
			Return(nil, testError)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(repoTxs)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		server := get_transactions.New(log, db).WithDependencies(deps)
		_, err := server.GetTransactions(
			ctx,
			testUserID,
			testLimit,
			testOffset,
			get_transactions.Amount,
			get_transactions.Desc)
		assert.Error(t, err)
	})

}
