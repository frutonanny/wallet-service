package get_report_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"

	"github.com/frutonanny/wallet-service/internal/repositories/report"
	"github.com/frutonanny/wallet-service/internal/services/get_report"
	mock "github.com/frutonanny/wallet-service/internal/services/get_report/mock"
)

const (
	testPeriod         = "2022-11"
	testPublicEndpoint = "publicEndpoint"
)

var (
	testError    = errors.New("error")
	testServices []report.Service
)

func TestService_GetReport(t *testing.T) {
	var db *sql.DB

	t.Run("get report successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		reportRepo := mock.NewMockRepository(ctrl)
		reportRepo.EXPECT().GetReport(ctx, testPeriod).Return(testServices, nil)

		minioClient := mock.NewMockMinioClient(ctrl)
		minioClient.
			EXPECT().
			PutObject(
				ctx,
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).
			Return(minio.UploadInfo{}, nil)

		log := mock.NewMocklogger(ctrl)

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(reportRepo)

		service := get_report.New(log, db, minioClient, testPublicEndpoint).WithDependencies(deps)
		_, err := service.GetReport(ctx, testPeriod)
		assert.NoError(t, err)
	})

	t.Run("get report failed, get report in repo error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		reportRepo := mock.NewMockRepository(ctrl)
		reportRepo.EXPECT().GetReport(ctx, testPeriod).Return(testServices, testError)

		minioClient := mock.NewMockMinioClient(ctrl)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(reportRepo)

		service := get_report.New(log, db, minioClient, testPublicEndpoint).WithDependencies(deps)
		_, err := service.GetReport(ctx, testPeriod)
		assert.Error(t, err)
	})

	t.Run("get report failed, put object error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()

		reportRepo := mock.NewMockRepository(ctrl)
		reportRepo.EXPECT().GetReport(ctx, testPeriod).Return(testServices, nil)

		minioClient := mock.NewMockMinioClient(ctrl)
		minioClient.
			EXPECT().
			PutObject(
				ctx,
				gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			).
			Return(minio.UploadInfo{}, testError)

		log := mock.NewMocklogger(ctrl)
		log.EXPECT().Error(gomock.Any())

		deps := mock.NewMockdependencies(ctrl)
		deps.EXPECT().NewRepository(gomock.Any()).Return(reportRepo)

		service := get_report.New(log, db, minioClient, testPublicEndpoint).WithDependencies(deps)
		_, err := service.GetReport(ctx, testPeriod)
		assert.Error(t, err)
	})
}
