package write_off_test

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
	write_off "github.com/frutonanny/wallet-service/internal/services/write-off"
	mock_write_off "github.com/frutonanny/wallet-service/internal/services/write-off/mock"
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
	testDelta      = int64(0)
	testBalance    = int64(1_000)
	testFailed     = int64(0)
)

var (
	testError = errors.New("error")
)

func TestService_WriteOff(t *testing.T) {
	t.Run("write-off cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().WriteOff(ctx, testWalletID, testAmount, testDelta).Return(testBalance, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusWrittenOff).Return(testTxID, nil)

		reportRepo := mock_write_off.NewMockReportRepository(ctrl)
		reportRepo.EXPECT().AddRecord(ctx, testServiceID, testAmount, gomock.Any()).Return(nil)

		txRepo := mock_write_off.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeWriteOff, gomock.Any(), testAmount).
			Return(testTxID, nil)

		mock.ExpectCommit()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)
		deps.EXPECT().NewReportRepository(gomock.Any()).Return(reportRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		balance, err := service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		require.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("write-off cash failed, ErrWalletNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testFailed, repositories.ErrRepoWalletNotFound)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_write_off.NewMocklogger(ctrl)

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrWalletNotFound)
	})

	t.Run("write-off cash failed, exist wallet error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testFailed, testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash cash failed, ErrOrderNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().
			GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, repositories.ErrRepoOrderNotFound)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_write_off.NewMocklogger(ctrl)

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		require.Error(t, err)
		assert.ErrorIs(t, err, servicesErrors.ErrOrderNotFound)
	})

	t.Run("write-off cash cash failed, get order error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash cash failed, update order error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash cash failed, add order transaction error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusWrittenOff).Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash cash failed, write-off error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().WriteOff(ctx, testWalletID, testAmount, testDelta).Return(testFailed, testError)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusWrittenOff).Return(testTxID, nil)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash failed, add transaction error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().WriteOff(ctx, testWalletID, testAmount, testDelta).Return(testBalance, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusWrittenOff).Return(testTxID, nil)

		txRepo := mock_write_off.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeWriteOff, gomock.Any(), testAmount).
			Return(testTxID, testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})

	t.Run("write-off cash failed, add report error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		walletRepo := mock_write_off.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)
		walletRepo.EXPECT().WriteOff(ctx, testWalletID, testAmount, testDelta).Return(testBalance, nil)

		orderRepo := mock_write_off.NewMockOrderRepository(ctrl)
		orderRepo.
			EXPECT().GetOrderByServiceID(ctx, testExternalID, testServiceID).
			Return(testOrderID, orders.StatusReserved, testAmount, nil)
		orderRepo.EXPECT().UpdateOrder(ctx, testOrderID, testAmount, orders.StatusWrittenOff).Return(nil)
		orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, orders.StatusWrittenOff).Return(testTxID, nil)

		txRepo := mock_write_off.NewMockTransactionRepository(ctrl)
		txRepo.EXPECT().AddTransaction(
			ctx, testWalletID, transactions.TypeWriteOff, gomock.Any(), testAmount).
			Return(testTxID, nil)

		reportRepo := mock_write_off.NewMockReportRepository(ctrl)
		reportRepo.EXPECT().AddRecord(ctx, testServiceID, testAmount, gomock.Any()).Return(testError)

		mock.ExpectRollback()

		deps := mock_write_off.NewMockdependencies(ctrl)
		deps.EXPECT().NewWalletRepository(gomock.Any()).Return(walletRepo)
		deps.EXPECT().NewOrderRepository(gomock.Any()).Return(orderRepo)
		deps.EXPECT().NewTransactionRepository(gomock.Any()).Return(txRepo)
		deps.EXPECT().NewReportRepository(gomock.Any()).Return(reportRepo)

		log := mock_write_off.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := write_off.New(log, db).WithDependencies(deps)

		_, err = service.WriteOff(ctx, testUserID, testServiceID, testExternalID, testAmount)
		assert.Error(t, err)
	})
}
