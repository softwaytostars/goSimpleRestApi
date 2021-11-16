package database

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DocumentCollectionName = "document"

type MongoDatastore struct {
	Database *mongo.Database
	Session  *mongo.Client
}

func (ds *MongoDatastore) Close() {
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
