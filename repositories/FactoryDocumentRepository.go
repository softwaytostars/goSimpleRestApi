package repositories

import (
	"goapi/config"
	"goapi/database"
)

func CreateDocumentRepository(config *config.Config) DocumentRepository {
	if config.StorageInMemory {
		return &InMemoryDocumentRepo{}
	} else {
		return &MongoDbDocumentRepo{database.GetMongoDataStore(config.DbConfig)}
	}
}
