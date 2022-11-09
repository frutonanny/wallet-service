package get_balance_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/repositories"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/get_balance"
	mock "github.com/frutonanny/wallet-service/internal/services/get_balance/mock"
)

const (
	testUserID   = int64(1)
	testFailed   = int64(0)
	testWalletID = int64(1)
	testBalance  = int64(10_000)
)

var testError = errors.New("error")

func TestService_GetBalance(t *testing.T) {
	var db *sql.DB

	t.Run("get balance successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := mock.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		repo.EXPECT().GetBalance(ctx, testWalletID).Return(testBalance, nil)

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)

		service := get_balance.New(log, db).WithDependencies(deps)

		balance, err := service.GetBalance(ctx, testUserID)
		require.NoError(t, err)
		assert.EqualValues(t, testBalance, balance)
	})

	t.Run("get balance failed, wallet not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := mock.NewMockRepository(ctrl)
		repo.
			EXPECT().
			ExistWallet(ctx, testUserID).
			Return(testFailed, repositories.ErrRepoWalletNotFound)

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_balance.New(log, db).WithDependencies(deps)

		_, err := service.GetBalance(ctx, testUserID)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("get balance failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		repo := mock.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		repo.EXPECT().GetBalance(ctx, testWalletID).Return(testBalance, testError)

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_balance.New(log, db).WithDependencies(deps)

		_, err := service.GetBalance(ctx, testUserID)
		require.Error(t, err)
	})
}
