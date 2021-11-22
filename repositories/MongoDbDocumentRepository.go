package repositories

import (
	"context"
	"errors"
	"goapi/database"
	"goapi/models"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDbDocumentRepo struct {
	store *database.MongoDatastore
}

//interface ObserverDatabase implementation
func (r *mongoDbDocumentRepo) SetDataStore(dataStore *database.MongoDatastore) {
	r.store = dataStore
}

func NewMongoDbDocumentRepo(databaseHandler *database.MongoDataBaseHandler) *mongoDbDocumentRepo {
	repo := &mongoDbDocumentRepo{}
	repo.store = databaseHandler.GetDataStore()
	if repo.store == nil {
		databaseHandler.RegisterAsObserver(repo)
	}
	return repo
}

func (r *mongoDbDocumentRepo) GetById(id string) (models.Document, error) {
	if r.store == nil {
		log.Error("data store not available")
		return models.Document{}, errors.New("no datastore")
	}

	//create collection
	collection := r.store.Database.Collection(database.DocumentCollectionName)

	//define filter
	filter := bson.D{primitive.E{Key: "id", Value: id}}

	var result models.Document
	//search
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Info("record does not exist")
		return models.Document{}, nil
	} else if err != nil {
		log.Error(err)
		return models.Document{}, err
	}
	return result, nil
}

func (r *mongoDbDocumentRepo) GetAll() ([]models.Document, error) {
	if r.store == nil {
		log.Error("data store not available")
		return nil, errors.New("no datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := r.store.Database.Collection(database.DocumentCollectionName)

	findOptions := options.Find()
	// Sort by `id` field ascending
	findOptions.SetSort(bson.D{primitive.E{Key: "id", Value: 1}})

	cur, err := collection.Find(ctx, bson.D{{}}, findOptions)

	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctxt *context.Context) {
		err := cur.Close(ctx)
		if err != nil {
			log.Error("Cannot close context", err)
		}
	}(cur, &ctx)

	if err := cur.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	var results []models.Document
	for cur.Next(ctx) {
		var result models.Document
		err := cur.Decode(&result)
		if err != nil {
			log.Error(err)
		} else {
			results = append(results, result)
		}
	}
	return results, nil
}

func (r *mongoDbDocumentRepo) CreateOrUpdate(document models.Document) (bool, error) {
	if r.store == nil {
		log.Error("data store not available")
		return false, errors.New("no datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//create collection
	collection := r.store.Database.Collection(database.DocumentCollectionName)

	//insert or update data
	filter := bson.M{"id": document.ID}

	pByte, err := bson.Marshal(document)
	if err != nil {
		logrus.Errorf("can't marshal:%s", err)
	}

	var update bson.M
	err = bson.Unmarshal(pByte, &update)
	if err != nil {
		logrus.Errorf("can't unmarshal:%s", err)
	}

	// NOTE: filter and ctx(Context) should be already defined
	res, err := collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})

	if err != nil {
		logrus.Error(err.Error())
		return false, err
	}

	if res.MatchedCount > 0 {
		return true, err
	}

	_, err = collection.InsertOne(ctx, document)

	return false, err
}

func (r *mongoDbDocumentRepo) Delete(id string) (bool, error) {
	if r.store == nil {
		log.Error("data store not available")
		return false, errors.New("no datastore")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := r.store.Database.Collection(database.DocumentCollectionName)

	//Define filter query for fetching specific document from collection
	filter := bson.D{primitive.E{Key: "id", Value: id}}

	result, err := collection.DeleteOne(ctx, filter)
	return result.DeletedCount > 0, err
}
