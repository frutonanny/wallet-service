package get_transactions_by_time_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/repositories"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time"
	mock_get_txs "github.com/frutonanny/wallet-service/internal/services/get_transactions_by_time/mock"
)

const (
	testUserID   = int64(1)
	testWalletID = int64(1)
	start        = "2022-07-02T10:00:00Z"
	end          = "2022-07-04T10:00:00Z"
)

var (
	testError    = errors.New("error")
	testStart, _ = time.Parse(time.RFC3339, start)
	testEnd, _   = time.Parse(time.RFC3339, end)
)

func TestService_GetTransactionsByTime(t *testing.T) {
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
			GetTransactionsByTime(ctx, testWalletID, testStart, testEnd).
			Return(nil, nil)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(repoTxs)

		log := mock_get_txs.NewMocklogger(ctrl)

		service := get_transactions_by_time.New(log, db).WithDependencies(deps)

		_, err := service.GetTransactionsByTime(ctx, testWalletID, testStart, testEnd)
		assert.NoError(t, err)
	})

	t.Run("get transaction, ErrWalletNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, repositories.ErrRepoWalletNotFound)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_transactions_by_time.New(log, db).WithDependencies(deps)

		_, err := service.GetTransactionsByTime(ctx, testWalletID, testStart, testEnd)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("get transaction, exist wallet error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, testError)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_transactions_by_time.New(log, db).WithDependencies(deps)

		_, err := service.GetTransactionsByTime(ctx, testWalletID, testStart, testEnd)
		assert.Error(t, err)
	})

	t.Run("get transaction, exist wallet error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repoWallet := mock_get_txs.NewMockWalletRepository(ctrl)
		repoWallet.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		repoTxs := mock_get_txs.NewMockTransactionRepository(ctrl)
		repoTxs.
			EXPECT().
			GetTransactionsByTime(ctx, testWalletID, testStart, testEnd).
			Return(nil, testError)

		deps := mock_get_txs.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(repoWallet)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(repoTxs)

		log := mock_get_txs.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_transactions_by_time.New(log, db).WithDependencies(deps)

		_, err := service.GetTransactionsByTime(ctx, testWalletID, testStart, testEnd)
		assert.Error(t, err)
	})
}
