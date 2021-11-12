package repositories

import (
	"goapi/config"
	"goapi/database"
)

func CloseRepositories(config *config.Config) {
	if !config.StorageInMemory {
		database.GetMongoDataStore(config.DbConfig).CloseMongoDataStore()
	}
}
