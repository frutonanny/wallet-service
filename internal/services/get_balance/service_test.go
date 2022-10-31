package get_balance_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/frutonanny/wallet-service/internal/repositories"
	internalErrors "github.com/frutonanny/wallet-service/internal/services/errors"
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

func TestGetBalance(t *testing.T) {
	var db *sql.DB

	t.Run("get balance successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(context.Background(), testUserID).Return(testWalletID, nil)
		repo.EXPECT().GetBalance(context.Background(), testWalletID).Return(testBalance, nil)

		builder := mock.NewMockBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := get_balance.New(log, db).WithBuilder(builder)

		balance, err := service.GetBalance(context.Background(), testUserID)
		assert.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("get balance failed, wallet not exist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		repo.
			EXPECT().
			ExistWallet(context.Background(), testUserID).
			Return(testFailed, repositories.ErrRepoWalletNotFound)

		builder := mock.NewMockBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_balance.New(log, db).WithBuilder(builder)

		_, err := service.GetBalance(context.Background(), testUserID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, internalErrors.ErrWalletNotFound)
	})

	t.Run("get balance failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		repo := mock.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(context.Background(), testUserID).Return(testWalletID, nil)
		repo.EXPECT().GetBalance(context.Background(), testWalletID).Return(testBalance, testError)

		builder := mock.NewMockBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := get_balance.New(log, db).WithBuilder(builder)

		_, err := service.GetBalance(context.Background(), testUserID)
		assert.Error(t, err)
		assert.ErrorIs(t, err, testError)
	})
}