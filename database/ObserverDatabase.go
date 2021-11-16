package database

type ObserverDatabase interface {
	SetDataStore(dataStore *MongoDatastore)
}
