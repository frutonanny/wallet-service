package cancel_test

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
	"github.com/frutonanny/wallet-service/internal/services/cancel"
	mock_cancel "github.com/frutonanny/wallet-service/internal/services/cancel/mock"
	servicesErrors "github.com/frutonanny/wallet-service/internal/services/errors"
	"github.com/frutonanny/wallet-service/internal/transactions"
)

const (
	testUserID     = int64(1)
	testWalletID   = int64(1)
	testOrderID    = int64(1)
	testTxID       = int64(0)
	testExternalID = int64(1)
	testAmount     = int64(1_000)
	testBalance    = int64(1_000)
	testFailed     = int64(0)
)

var (
	testError = errors.New("error")
)

func TestService_Cancel(t *testing.T) {
	t.Run("cancel reservation cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Cancel(ctx, testWalletID, testAmount).Return(testBalance, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, orders.StatusCancelled).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusCancelled).Return(testTxID, nil)

		txRepo := mock_cancel.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeCancel, gomock.Any(), testAmount).
			Return(testTxID, nil)

		mock.ExpectCommit()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		balance, err := service.Cancel(ctx, testUserID, testExternalID)
		require.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("cancel reservation failed, ErrWalletNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testFailed, repositories.ErrRepoWalletNotFound)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_cancel.NewMocklogger(ctrl)

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("cancel reservation failed, exist wallet error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testFailed, testError)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		require.Error(t, err)
	})

	t.Run("cancel reservation cash failed, ErrOrderNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			GetOrder(ctx, testExternalID).
			Return(testOrderID, orders.StatusReserved, testAmount, repositories.ErrRepoOrderNotFound)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_cancel.NewMocklogger(ctrl)

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrOrderNotFound)
	})

	t.Run("cancel reservation cash failed, get order error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			GetOrder(ctx, testExternalID).
			Return(testOrderID, orders.StatusReserved, testAmount, testError)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		assert.Error(t, err)
	})

	t.Run("cancel reservation cash failed, update order error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, orders.StatusCancelled).Return(testError)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		assert.Error(t, err)
	})

	t.Run("cancel reservation cash failed, add order transaction error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, orders.StatusCancelled).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusCancelled).Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		assert.Error(t, err)
	})

	t.Run("cancel reservation cash failed, cancel error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Cancel(ctx, testWalletID, testAmount).Return(testFailed, testError)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, orders.StatusCancelled).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusCancelled).Return(testTxID, nil)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		assert.Error(t, err)
	})

	t.Run("cancel reservation cash failed, add transaction error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().Cancel(ctx, testWalletID, testAmount).Return(testBalance, nil)

		orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
		orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, orders.StatusCancelled).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusCancelled).Return(testTxID, nil)

		txRepo := mock_cancel.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeCancel, gomock.Any(), testAmount).
			Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_cancel.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)

		log := mock_cancel.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := cancel.New(log, db).WithDependencies(deps)

		_, err = service.Cancel(ctx, testUserID, testExternalID)
		assert.Error(t, err)
	})
}
