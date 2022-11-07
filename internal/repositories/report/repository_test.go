package report_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	repoReport "github.com/frutonanny/wallet-service/internal/repositories/report"
	testingboilerplate "github.com/frutonanny/wallet-service/internal/testing_boilerplate"
)

const (
	fileConfig    = "../../../config/config.local.json"
	testServiceID = int64(1)
	testAmount    = int64(100)
)

var (
	config = serviceConfig.Must(fileConfig)
)

func TestRepository_AddRecord(t *testing.T) {
	t.Run("add record", func(t *testing.T) {
		ctx := context.Background()

		period := time.Now()
		periodFormatted := period.Format(repoReport.PeriodLayout)

		tx, cancel := testingboilerplate.InitDB(t, config.DB.DSN)
		defer cancel()

		repo := repoReport.New(tx)

		err := repo.AddRecord(ctx, testServiceID, testAmount, period)
		require.NoError(t, err)

		report, err := repo.GetReport(ctx, periodFormatted)
		require.NoError(t, err)
		require.Len(t, report, 1)
		assert.Equal(t, testServiceID, report[0].ServiceID)
		assert.EqualValues(t, testAmount, report[0].TotalRevenue)

		err = repo.AddRecord(ctx, testServiceID, testAmount, period)
		require.NoError(t, err)

		report, err = repo.GetReport(ctx, periodFormatted)
		require.NoError(t, err)
		require.Len(t, report, 1)
		assert.Equal(t, testServiceID, report[0].ServiceID)
		assert.EqualValues(t, testAmount*2, report[0].TotalRevenue)
	})
}
