package add_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/frutonanny/wallet-service/internal/repositories"
	"github.com/frutonanny/wallet-service/internal/services/add"
	mock_add "github.com/frutonanny/wallet-service/internal/services/add/mock"
)

const (
	testUserID   = int64(1)
	testWalletID = int64(1)
	testCash     = int64(1_000)
	testBalance  = int64(1_000)
)

var testError = errors.New("error")

func TestAdd(t *testing.T) {
	t.Run("add cash successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		repo := mock_add.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(context.Background(), testUserID).Return(testWalletID, nil)
		repo.EXPECT().Add(context.Background(), testWalletID, testCash).Return(testBalance, nil)
		repo.EXPECT().
			AddTransaction(context.Background(), testWalletID, gomock.Any(), gomock.Any(), testCash).
			Return(nil)

		mock.ExpectCommit()

		builder := mock_add.NewMockRepoBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := add.New(log, db).WithBuilder(builder)

		balance, err := service.Add(context.Background(), testUserID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("add cash successfully, created wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		repo := mock_add.NewMockRepository(ctrl)
		repo.EXPECT().
			ExistWallet(context.Background(), testUserID).
			Return(int64(0), repositories.ErrRepoWalletNotFound)
		repo.EXPECT().CreateWallet(context.Background(), testUserID).Return(testWalletID, nil)
		repo.EXPECT().Add(context.Background(), testWalletID, testCash).Return(testBalance, nil)
		repo.EXPECT().
			AddTransaction(context.Background(), testWalletID, gomock.Any(), gomock.Any(), testCash).
			Return(nil)

		mock.ExpectCommit()

		builder := mock_add.NewMockRepoBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Info(gomock.Any())

		service := add.New(log, db).WithBuilder(builder)

		balance, err := service.Add(context.Background(), testUserID, testCash)
		assert.NoError(t, err)
		assert.Equal(t, testBalance, balance)
	})

	t.Run("add cash failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		db, mock, err := sqlmock.New()
		require.NoError(t, err)

		mock.ExpectBegin()

		repo := mock_add.NewMockRepository(ctrl)
		repo.EXPECT().ExistWallet(context.Background(), testUserID).Return(testWalletID, testError)

		mock.ExpectRollback()

		builder := mock_add.NewMockRepoBuilder(ctrl)
		builder.EXPECT().NewRepository(gomock.Any()).Return(repo)

		log := mock_add.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		service := add.New(log, db).WithBuilder(builder)

		_, err = service.Add(context.Background(), testUserID, testCash)
		assert.Error(t, err)
	})
}
