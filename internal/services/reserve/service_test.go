package reserve_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/orders"
	"github.com/frutonanny/wallet-service/internal/repositories"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/services/reserve"
	mock_reserve "github.com/frutonanny/wallet-service/internal/services/reserve/mock"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

const (
	testUserID     = int64(1)
	testWalletID   = int64(1)
	testOrderID    = int64(1)
	testTxID       = int64(0)
	testExternalID = int64(1)
	testServiceID  = int64(1)
	testAmount     = int64(1_000)
	testBalance    = int64(1_000)
)

var testError = errors.New("error")

func TestService_Reserve(t *testing.T) {
	t.Run("reservation cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Reserve(ctx, testWalletID, testAmount).Return(testBalance, nil)

		orderRepo := mock_reserve.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			CreateOrder(ctx, testWalletID, testExternalID, testServiceID, testAmount, orders.StatusReserved).
			Return(testOrderID, nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusReserved).Return(testTxID, nil)

		txRepo := mock_reserve.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeReserve, gomock.Any(), testAmount).
			Return(testTxID, nil)

		mock.ExpectCommit()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)

		log := mock_reserve.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := reserve.New(log, db).WithDependencies(deps)

		balance, err := service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		require.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("reservation cash failed, ErrWalletNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, repositories.ErrRepoWalletNotFound)

		mock.ExpectRollback()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_reserve.NewMocklogger(ctrl)

		service := reserve.New(log, db).WithDependencies(deps)

		_, err = service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("reservation cash failed, wallet exist error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, testError)

		mock.ExpectRollback()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_reserve.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := reserve.New(log, db).WithDependencies(deps)

		_, err = service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("reservation cash failed, reserved error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Reserve(ctx, testWalletID, testAmount).Return(testBalance, testError)

		mock.ExpectRollback()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_reserve.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := reserve.New(log, db).WithDependencies(deps)

		_, err = service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("reservation cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Reserve(ctx, testWalletID, testAmount).Return(testBalance, nil)

		orderRepo := mock_reserve.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			CreateOrder(ctx, testWalletID, testExternalID, testServiceID, testAmount, orders.StatusReserved).
			Return(testOrderID, nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusReserved).Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_reserve.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := reserve.New(log, db).WithDependencies(deps)

		_, err = service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("reservation cash failed, add transaction error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_reserve.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Reserve(ctx, testWalletID, testAmount).Return(testBalance, nil)

		orderRepo := mock_reserve.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			CreateOrder(ctx, testWalletID, testExternalID, testServiceID, testAmount, orders.StatusReserved).
			Return(testOrderID, nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusReserved).Return(testTxID, nil)

		txRepo := mock_reserve.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeReserve, gomock.Any(), testAmount).
			Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_reserve.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)

		log := mock_reserve.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := reserve.New(log, db).WithDependencies(deps)

		_, err = service.Reserve(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})
}
