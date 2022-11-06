package get_transactions_by_time

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoTxs "github.com/frutonanny/wallet-service/internal/repositories/transaction"
	repoWallet "github.com/frutonanny/wallet-service/internal/repositories/wallet"
)

type dependenciesImpl struct{}

func (b *dependenciesImpl) NewWalletRepository(db postgres.Database) walletRepository {
	return repoWallet.New(db)
}

func (b *dependenciesImpl) NewTransactionRepository(db postgres.Database) transactionRepository {
	return repoTxs.New(db)
}
