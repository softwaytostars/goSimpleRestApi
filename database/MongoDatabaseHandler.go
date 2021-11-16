package database

import (
	"context"
	"goapi/config"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDataBaseHandler struct {
	observers []ObserverDatabase
	dataStore *MongoDatastore
	lock      sync.RWMutex
}

var instanceDBHandler *MongoDataBaseHandler
var onceDBHandler sync.Once

func (h *MongoDataBaseHandler) RegisterAsObserver(o ObserverDatabase) {
	h.observers = append(h.observers, o)
}

func (h *MongoDataBaseHandler) GetDataStore() *MongoDatastore {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.dataStore
}

func (h *MongoDataBaseHandler) notifyObserversDataStoreIsAvailable() {
	for _, o := range h.observers {
		o.SetDataStore(h.dataStore)
	}
}

func (h *MongoDataBaseHandler) TryOrRetryCreateConnection(config *config.DatabaseConfig) {
	//schedule try connection every 10 seconds util success
	s := gocron.NewScheduler(time.UTC)
	s.Every(10).Seconds().Do(func() {
		log.Info("Attempt to connect to DB")
		db, session, err := connectToMongo(config)
		if err == nil {
			log.Info("Successfull connection")
			datastore := &MongoDatastore{db, session}
			datastore.createDatabaseForApp()
			h.lock.Lock()
			defer h.lock.Unlock()
			h.dataStore = datastore
			h.notifyObserversDataStoreIsAvailable()
			s.Stop() //stop the cron
		}
	})
	s.StartAsync()
}

func connectToMongo(dbConfig *config.DatabaseConfig) (a *mongo.Database, b *mongo.Client, err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
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

func (h *MongoDataBaseHandler) Close() {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.dataStore.Close()
}

func GetMongoDatabaseHandler() *MongoDataBaseHandler {
	onceDBHandler.Do(func() {
		instanceDBHandler = &MongoDataBaseHandler{}
	})
	return instanceDBHandler
}
