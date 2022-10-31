package add

import (
	"github.com/frutonanny/wallet-service/internal/postgres"
	repositoryWallet "github.com/frutonanny/wallet-service/internal/repositories/wallet"
)

type builderImpl struct{}

func (b *builderImpl) NewRepository(db postgres.Database) Repository {
	return repositoryWallet.New(db)
}
