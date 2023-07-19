package repository

import "github.com/ishanshre/GoRestAPIMongoDB/internals/models"

type MongoDbRepo interface {
	GetAllUsers(page, limit int) ([]*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	DeleteUser(username string) error
	UsernameExists(username string) error
}
