package get_balance

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repositoryWallet "github.com/frutonanny/wallet-service/internal/repositories/wallet"
)

type dependenciesImpl struct{}

func (b *dependenciesImpl) NewRepository(db postgres.Database) Repository {
	return repositoryWallet.New(db)
}
