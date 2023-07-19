package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbInterface interface {
	GetUserCollection() *mongo.Collection
}

type DB struct {
	Client *mongo.Client
}

func ConnectToNoSql(dsn string) (DbInterface, error) {
	client, err := NewDatabase(dsn)
	if err != nil {
		return nil, err
	}
	return &DB{
		Client: client,
	}, nil
}

func (db *DB) GetUserCollection() *mongo.Collection {
	return db.Client.Database("myDB").Collection("users")
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
