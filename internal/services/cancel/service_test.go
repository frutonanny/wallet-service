package cancel

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/frutonanny/wallet-service/internal/orders"
	mock_cancel "github.com/frutonanny/wallet-service/internal/services/cancel/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	testUserID     = int64(1)
	testWalletID   = int64(1)
	testOrderID    = int64(1)
	testExternalID = int64(1)
	testServiceID  = int64(1)
	testAmount     = int64(1_000)
	testBalance    = int64(1_000)
)

var (
	testError  = errors.New("error")
	testStatus = orders.StatusReserved
)

func TestCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectBegin()

	walletRepo := mock_cancel.NewMockWalletRepository(ctrl)
	walletRepo.EXPECT().ExistWallet(ctx, testUserID).Return(testWalletID, nil)

	orderRepo := mock_cancel.NewMockOrderRepository(ctrl)
	orderRepo.EXPECT().GetOrder(ctx, testExternalID).Return(testOrderID, testStatus, testAmount, nil)
	orderRepo.EXPECT().UpdateOrderStatus(ctx, testOrderID, testStatus).Return(nil)
	orderRepo.EXPECT().AddOrderTransactions(ctx, testOrderID, testStatus).Return(nil)

	txRepo := mock_cancel.NewMockTransactionRepository(ctrl)
	txRepo.EXPECT().AddTransaction(ctx, testWalletID, gomock.Any(), gomock.Any(), testAmount).Return(nil)

	mock.ExpectCommit()

	deps := mock_cancel.NewMockdependencies(ctrl)
	deps.EXPECT().NewWalletRepository(ctrl)
	deps.EXPECT().NewOrderRepository(ctrl)
	deps.EXPECT().NewTransactionRepository(ctrl)

	log := mock_cancel.NewMocklogger(ctrl)
	log.EXPECT().Info(gomock.Any())

	//service := New(log, db).WithDependencies(deps)

}
