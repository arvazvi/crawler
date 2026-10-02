package db

import (
	"context"

	"github.com/arvazvi/crawler/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client *mongo.Client
	collection *mongo.Collection
}

func NewConnection(ctx context.Context, conf *config.Config) (*MongoDB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DB.URI))
	if err != nil {
		return nil, err
	}

	coll := client.Database(conf.DB.Database).Collection(conf.DB.Collection)
	filter := bson.D{{}}
	coll.DeleteMany(ctx, filter)

	return &MongoDB{client: client, collection: coll}, nil
}

func (d *MongoDB) Close(ctx context.Context) {
	d.client.Disconnect(ctx)
}

func (d *MongoDB) Insert(ctx context.Context, object any) error {
	_, err := d.collection.InsertOne(context.TODO(), object)
	return err
}