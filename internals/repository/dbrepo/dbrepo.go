package dbrepo

import (
	"context"
	"time"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/database"
	"github.com/ishanshre/GoRestAPIMongoDB/internals/repository"
)

type mongodbRepo struct {
	Client database.DbInterface
	ctx    context.Context
}

func NewMongoDbRepo(client database.DbInterface, ctx context.Context) repository.MongoDbRepo {
	return &mongodbRepo{
		Client: client,
		ctx:    ctx,
	}
}

const timeout = 3 * time.Second
