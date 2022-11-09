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

		// Добавляем данные в отчет за период period.
		err := repo.AddRecord(ctx, testServiceID, testAmount, period)
		require.NoError(t, err)

		// Получаем данные из отчета и проверяем, что они соответствуют ранее занесенным.
		report, err := repo.GetReport(ctx, periodFormatted)
		require.NoError(t, err)
		require.Len(t, report, 1)
		assert.Equal(t, testServiceID, report[0].ServiceID)
		assert.EqualValues(t, testAmount, report[0].TotalRevenue)

		// Обновляем данные в отчете для того же сервиса и за тот же период period.
		// Ожидаем, что запись обновится, увеличится TotalRevenue на сумму testAmount.
		err = repo.AddRecord(ctx, testServiceID, testAmount, period)
		require.NoError(t, err)

		// Получаем данные из отчета и проверяем, что они обновились. TotalRevenue  = 2*testAmount.
		report, err = repo.GetReport(ctx, periodFormatted)
		require.NoError(t, err)
		require.Len(t, report, 1)
		assert.Equal(t, testServiceID, report[0].ServiceID)
		assert.EqualValues(t, testAmount*2, report[0].TotalRevenue)
	})
}
