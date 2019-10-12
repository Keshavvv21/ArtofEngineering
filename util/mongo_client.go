package util

import (
	"context"

	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate counterfeiter -o fakes/fake_mongo_client.go --fake-name FakeMongoClient . MongoClient
type MongoClient interface {
	GetCollection(dbName, collectionName string) *mongo.Collection
	InsertData(dbName, collectionName string, data interface{}) (*mongo.InsertOneResult, error)
	FindObject(dbName,
		collectionName string, filter, result interface{}) error
	GetDatabase(dbName string) *mongo.Database
	FindAllObjects(dbName,
		collectionName string, result interface{}, limit int64) error
	UpdateOne(dbName, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error)
	FindObjects(dbName,
		collectionName string, filter, result interface{}) error
}

type MongoClientImpl struct {
	MongoClient *mongo.Client
	ctx         context.Context
}

func NewMongoClient(ctx context.Context, url string) MongoClient {
	client, _ := CreateClient(ctx, url)
	return &MongoClientImpl{
		client,
		ctx,
	}
}

func (mg *MongoClientImpl) GetCollection(dbName, collectionName string) *mongo.Collection {
	glog.Info("get-collection-started")
	defer glog.Info("get-collection-completed")
	collection := mg.GetDatabase(dbName).Collection(collectionName)
	return collection
}

func (mg *MongoClientImpl) GetDatabase(dbName string) *mongo.Database {
	glog.Info("get-database-started")
	defer glog.Info("get-database-completed")
	database := mg.MongoClient.Database(dbName)
	return database
}

func (mg *MongoClientImpl) InsertData(dbName, collectionName string, data interface{}) (*mongo.InsertOneResult, error) {
	glog.Info("insert-data-started")
	defer glog.Info("insert-data-completed")
	collection := mg.GetCollection(dbName, collectionName)
	result, err := collection.InsertOne(mg.ctx, data)
	return result, err
}

func (mg *MongoClientImpl) UpdateOne(dbName, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	glog.Info("update-data-started")
	defer glog.Info("update-data-completed")
	collection := mg.GetCollection(dbName, collectionName)
	result, err := collection.UpdateOne(mg.ctx, filter, update)
	return result, err
}

func (mg *MongoClientImpl) FindObject(dbName,
	collectionName string, filter, result interface{}) error {
	glog.Info("find-object-started")
	defer glog.Info("find-object-completed")
	collection := mg.GetCollection(dbName, collectionName)
	err := collection.FindOne(mg.ctx, filter).Decode(&result)
	return err
}

func (mg *MongoClientImpl) FindObjects(dbName,
	collectionName string, filter, result interface{}) error {
	glog.Info("find-objects-started")
	defer glog.Info("find-objects-completed")
	collection := mg.GetCollection(dbName, collectionName)
	cursors, err := collection.Find(mg.ctx, filter)
	if err != nil {
		return nil
	}
	cursors.Decode(result)
	return nil
}

func (mg *MongoClientImpl) FindAllObjects(dbName,
	collectionName string, result interface{}, limit int64) error {
	glog.Info("find-all-objects-started")
	defer glog.Info("find-all-objects-completed")
	findOptions := options.Find()
	findOptions.SetLimit(limit)
	collection := mg.GetCollection(dbName, collectionName)
	cursors, err := collection.Find(mg.ctx, bson.D{{}}, findOptions)
	if err != nil {
		return nil
	}
	cursors.Decode(result)
	return nil
}

func CreateClient(ctx context.Context, url string) (*mongo.Client, error) {
	glog.Info("creating-monogo-client-started")
	defer glog.Info("creating-monogo-client-completed")
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		glog.Error("check-mongo-connection", err)
		return nil, err
	}
	return client, nil
}
