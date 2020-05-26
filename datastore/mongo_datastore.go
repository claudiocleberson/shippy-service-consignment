package datastore

import (
	"context"
	"log"
	"time"

	"github.com/claudiocleberson/shippy-service-consignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dabaseName     = "shippy"
	collectionName = "consignments"
)

var (
	mongoCollection *mongo.Collection
	retry           int
)

type MongoClient interface {
	Create(context.Context, *models.Consignment) error
	GetAll(context.Context) ([]*models.Consignment, error)
}

func NewMongoClient(uri string) MongoClient {
	connectMongoCluster(uri)

	return &mongoClient{}
}

func connectMongoCluster(uri string) {

	log.Println("Connecting mongo cluster...")

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {

		if retry >= 3 {
			panic(err)
			return
		}
		retry = retry + 1
		time.Sleep(time.Second * 2)
		connectMongoCluster(uri)
	}

	mongoCollection = client.Database(dabaseName).Collection(collectionName)

	log.Println("Mongo cluster connected...")
}

type mongoClient struct{}

func (m *mongoClient) Create(ctx context.Context, cons *models.Consignment) error {
	_, err := mongoCollection.InsertOne(ctx, cons)
	return err
}

func (m *mongoClient) GetAll(ctx context.Context) ([]*models.Consignment, error) {

	cur, err := mongoCollection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	var consignments []*models.Consignment
	for cur.Next(ctx) {
		var consignment models.Consignment
		if err := cur.Decode(&consignment); err != nil {
			return nil, err
		}
		consignments = append(consignments, &consignment)
	}

	return consignments, nil
}
