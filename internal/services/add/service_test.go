package add_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/services/add"
	mock_add "github.com/frutonanny/wallet-service/internal/services/add/mock"
)

const (
	testUserID   = int64(1)
	testWalletID = int64(1)
	testAmount   = int64(1_000)
	testBalance  = int64(1_000)
	testTxID     = int64(0)
)

var testError = errors.New("error")

func TestAdd(t *testing.T) {
	t.Run("add cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_add.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateIfNotExist(context.Background(), testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Add(context.Background(), testWalletID, testAmount).Return(testBalance, nil)

		txsRepo := mock_add.NewMockTransactionRepository(ctrl)
		txsRepo.EXPECT().
			AddTransaction(context.Background(), testWalletID, gomock.Any(), gomock.Any(), testAmount).
			Return(testTxID, nil)

		mock.ExpectCommit()

		deps := mock_add.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txsRepo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := add.New(log, db).WithDependencies(deps)

		balance, err := service.Add(context.Background(), testUserID, testAmount)
		require.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("add cash failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_add.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateIfNotExist(context.Background(), testUserID).Return(testWalletID, testError)

		mock.ExpectRollback()

		deps := mock_add.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := add.New(log, db).WithDependencies(deps)

		_, err = service.Add(context.Background(), testUserID, testAmount)
		assert.Error(t, err)
	})

	t.Run("add cash failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_add.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateIfNotExist(context.Background(), testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Add(context.Background(), testWalletID, testAmount).Return(testBalance, testError)

		mock.ExpectRollback()

		deps := mock_add.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := add.New(log, db).WithDependencies(deps)

		_, err = service.Add(context.Background(), testUserID, testAmount)
		assert.Error(t, err)
	})

	t.Run("add cash failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_add.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateIfNotExist(context.Background(), testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Add(context.Background(), testWalletID, testAmount).Return(testBalance, nil)

		txsRepo := mock_add.NewMockTransactionRepository(ctrl)

		txsRepo.EXPECT().
			AddTransaction(context.Background(), testWalletID, gomock.Any(), gomock.Any(), testAmount).
			Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_add.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txsRepo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := add.New(log, db).WithDependencies(deps)

		_, err = service.Add(context.Background(), testUserID, testAmount)
		assert.Error(t, err)
	})
}
