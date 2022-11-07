package write_off

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoOrder "github.com/frutonanny/wallet-service/internal/repositories/order"
	repoReport "github.com/frutonanny/wallet-service/internal/repositories/report"
	repoTxs "github.com/frutonanny/wallet-service/internal/repositories/transaction"
	repoWallet "github.com/frutonanny/wallet-service/internal/repositories/wallet"
)

type dependenciesImpl struct{}

func (b *dependenciesImpl) NewWalletRepository(db postgres.Database) WalletRepository {
	return repoWallet.New(db)
}
func (b *dependenciesImpl) NewOrderRepository(db postgres.Database) OrderRepository {
	return repoOrder.New(db)
}

func (b *dependenciesImpl) NewTransactionRepository(db postgres.Database) TransactionRepository {
	return repoTxs.New(db)
}

func (b *dependenciesImpl) NewReportRepository(db postgres.Database) ReportRepository {
	return repoReport.New(db)
}
