package cancel

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repoOrder "github.com/frutonanny/wallet-service/internal/repositories/order"
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
