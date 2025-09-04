package util

import (
	"context"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//
// MongoClient Interface
//
// This defines all database operations that the application needs.
// By depending on the interface (instead of concrete implementation),
// we can easily mock the database layer for testing.
//
 
//go:generate counterfeiter -o fakes/fake_mongo_client.go --fake-name FakeMongoClient . MongoClient
type MongoClient interface {
	GetCollection(dbName, collectionName string) *mongo.Collection
	GetDatabase(dbName string) *mongo.Database
	InsertData(dbName, collectionName string, data interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(dbName, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error)
	FindObject(dbName, collectionName string, filter, result interface{}) error
	FindObjects(dbName, collectionName string, filter, result interface{}) error
	FindAllObjects(dbName, collectionName string, result interface{}, limit int64) error
}

//
// MongoClientImpl
//
// Concrete implementation of MongoClient interface.
// Holds a mongo.Client and a context that is reused across DB operations.
//
type MongoClientImpl struct {
	MongoClient *mongo.Client
	ctx         context.Context
}

//
// Factory method for creating a new MongoClient.
// Accepts context and connection URL, creates a new client, and returns wrapper.
//
func NewMongoClient(ctx context.Context, url string) MongoClient {
	client, _ := CreateClient(ctx, url)
	return &MongoClientImpl{
		MongoClient: client,
		ctx:         ctx,
	}
}

//
// GetCollection: returns a MongoDB collection reference.
//
func (mg *MongoClientImpl) GetCollection(dbName, collectionName string) *mongo.Collection {
	glog.Info("get-collection-started")
	defer glog.Info("get-collection-completed")

	return mg.GetDatabase(dbName).Collection(collectionName)
}

//
// GetDatabase: returns a MongoDB database reference.
//
func (mg *MongoClientImpl) GetDatabase(dbName string) *mongo.Database {
	glog.Info("get-database-started")
	defer glog.Info("get-database-completed")

	return mg.MongoClient.Database(dbName)
}

//
// InsertData: inserts a single document into a collection.
//
func (mg *MongoClientImpl) InsertData(dbName, collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	glog.Info("insert-data-started")
	defer glog.Info("insert-data-completed")

	collection := mg.GetCollection(dbName, collectionName)
	return collection.InsertOne(mg.ctx, data)
}

//
// UpdateOne: updates a single document matching filter.
//
func (mg *MongoClientImpl) UpdateOne(dbName, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	glog.Info("update-data-started")
	defer glog.Info("update-data-completed")

	collection := mg.GetCollection(dbName, collectionName)
	return collection.UpdateOne(mg.ctx, filter, update)
}

//
// FindObject: finds a single document matching filter and decodes into result.
//
func (mg *MongoClientImpl) FindObject(dbName, collectionName string, filter, result interface{}) error {
	glog.Info("find-object-started")
	defer glog.Info("find-object-completed")

	collection := mg.GetCollection(dbName, collectionName)
	return collection.FindOne(mg.ctx, filter).Decode(result)
}

//
// FindObjects: finds all documents matching filter and decodes into result slice.
//

func (mg *MongoClientImpl) FindObjects(dbName, collectionName string, filter, result interface{}) error {
	glog.Info("find-objects-started")
	defer glog.Info("find-objects-completed")

	collection := mg.GetCollection(dbName, collectionName)
	cursor, err := collection.Find(mg.ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(mg.ctx)

	return cursor.All(mg.ctx, result) // safer than Decode
}

//
// FindAllObjects: fetches all documents (with optional limit) and decodes into result slice.
//
func (mg *MongoClientImpl) FindAllObjects(dbName, collectionName string, result interface{}, limit int64) error {
	glog.Info("find-all-objects-started")
	defer glog.Info("find-all-objects-completed")

	findOptions := options.Find().SetLimit(limit)
	collection := mg.GetCollection(dbName, collectionName)

	cursor, err := collection.Find(mg.ctx, bson.D{}, findOptions)
	if err != nil {
		return err
	}
	defer cursor.Close(mg.ctx)

	return cursor.All(mg.ctx, result) // safer than Decode
}

//
// CreateClient: connects to MongoDB and verifies connection with Ping.
//
func CreateClient(ctx context.Context, url string) (*mongo.Client, error) {
	glog.Info("creating-mongo-client-started")
	defer glog.Info("creating-mongo-client-completed")

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection with a ping
	if err := client.Ping(ctx, nil); err != nil {
		glog.Error("mongo-connection-ping-failed", err)
		return nil, err
	}

	return client, nil
}
