package database

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goapi/config"
	"sync"
	"sync/atomic"
	"time"
)

const DocumentCollectionName = "document"

type MongoDatastore struct {
	Database *mongo.Database
	Session  *mongo.Client
}

type onceIfSuccess struct {
	done uint32
	m    sync.Mutex
}

func (o *onceIfSuccess) Do(f func() error) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

func (o *onceIfSuccess) doSlow(f func() error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		err := f()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
}

var instance *MongoDatastore
var once onceIfSuccess

func GetMongoDataStore(config config.DatabaseConfig) *MongoDatastore {
	once.Do(func() error {
		db, session, err := connectToMongo(config)
		if err == nil {
			instance = &MongoDatastore{db, session}
			instance.createDatabaseForApp()
		} else {
			log.Error("MongoDB: Failed to connect to database: %v", config.DBName)
		}
		return err
	})
	return instance
}

func connectToMongo(dbConfig config.DatabaseConfig) (a *mongo.Database, b *mongo.Client, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConfig.Uri).SetMaxPoolSize(dbConfig.MaxPoolSize))
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error(err)
		errDisco := client.Disconnect(ctx)
		if errDisco != nil {
			log.Error("Cannot close DB client")
		}
		return nil, nil, err
	}

	return client.Database(dbConfig.DBName), client, nil
}

func (ds *MongoDatastore) CloseMongoDataStore() {
	err := ds.Session.Disconnect(context.Background())
	if err != nil {
		log.Error("Cannot disconnect DB client")
	}
}

func (ds *MongoDatastore) createDatabaseForApp() {
	//create model index
	mod := mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}

	//create collection
	collection := ds.Database.Collection(DocumentCollectionName)

	//create context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//create index
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Errorf("Cannot create index on %s", DocumentCollectionName)
	}
}
