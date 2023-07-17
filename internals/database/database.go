package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client         *mongo.Client
	Collectiontion *mongo.Collection
}

func ConnectToNoSql(dsn string) (*DB, error) {
	client, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}
	collection := CreateCollections(client, "auth", "users")
	return &DB{
		Client:         client,
		Collectiontion: collection,
	}, nil
}

func CreateCollections(c *mongo.Client, database, name string) *mongo.Collection {
	return c.Database(database).Collection(name)
}

func NewDatabase(dsn string) (*mongo.Client, error) {

	// create a context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// create options for client
	clientOptions := options.Client().ApplyURI(dsn)

	// connect to the mongo db
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// ping the database to check its connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return client, err
}
