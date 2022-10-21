package main

import (
	"fmt"

	serviceConfig "github.com/frutonanny/wallet-service/internal/config"
	serviceLogger "github.com/frutonanny/wallet-service/internal/logger"
	serviceDB "github.com/frutonanny/wallet-service/internal/postgres"
)

const path = "config/config.json"

func main() {
	config := serviceConfig.Must(path)

	logger := serviceLogger.Must()

	db := serviceDB.MustConnect(config.DB.DSN)
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(fmt.Sprintf("close db error: %s", err))
		}
	}()

	serviceDB.MustMigrate(db)
}
