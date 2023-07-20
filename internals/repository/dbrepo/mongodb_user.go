package dbrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ishanshre/GoRestAPIMongoDB/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllUser returns the slice of users and error.
// It takes limit and page as paramaters.
func (m *mongodbRepo) GetAllUsers(page, limit int) ([]*models.User, error) {
	// set limit timeout
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	// setup pagination condition
	if page == 0 || page < 0 {
		page = 1
	}
	if limit == 0 || limit < 0 {
		limit = 10
	}
	skip := (page - 1) * limit
	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	query := bson.M{}

	result, err := m.Client.GetUserCollection().Find(ctx, query, &opt)
	if err != nil {
		return nil, errors.New("error in fetching users")
	}
	defer result.Close(m.ctx)
	users := []*models.User{}
	for result.Next(m.ctx) {
		user := &models.User{}
		if err := result.Decode(&user); err != nil {
			return nil, errors.New("error in unmarshling users")
		}
		users = append(users, user)
	}
	return users, nil
}

func (m *mongodbRepo) CreateUser(user *models.User) (*models.User, error) {
	// config time out
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	user.CreatedAt = time.Now()

	// insert into mongo db
	res, err := m.Client.GetUserCollection().InsertOne(ctx, user)
	if err != nil {
		return nil, errors.New("error in inserting user")
	}
	user.ID = res.InsertedID.(primitive.ObjectID)

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: opt,
	}
	_, err = m.Client.GetUserCollection().Indexes().CreateOne(ctx, index)
	if err != nil {
		return nil, errors.New("error in creating unique index for user username")
	}
	newUser := &models.User{}
	query := bson.M{"_id": res.InsertedID}
	if err := m.Client.GetUserCollection().FindOne(ctx, query).Decode(&newUser); err != nil {
		return nil, errors.New("error in fetching newly created user")
	}
	return newUser, nil
}

func (m *mongodbRepo) DeleteUser(username string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	res, err := m.Client.GetUserCollection().DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		return errors.New("error in deleting user")
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("delete count: %d", res.DeletedCount)
	}
	return nil

}

func (m *mongodbRepo) UsernameExists(username string) error {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()
	existingUser := &models.User{}
	err := m.Client.GetUserCollection().FindOne(ctx, bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return errors.New("username already exists")
	}
	return nil
}

func (m *mongodbRepo) GetUserByUsername(username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()

	query := bson.M{"username": username}
	user := &models.User{}
	if err := m.Client.GetUserCollection().FindOne(ctx, query).Decode(&user); err != nil {
		return nil, fmt.Errorf("error in fetcthing the user with username %s", username)
	}
	return user, nil
}

func (m *mongodbRepo) UpdateUser(username string, updateObj *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(m.ctx, timeout)
	defer cancel()
	query := bson.M{"username": username}
	update := bson.D{{"$set", bson.D{
		{"first_name", updateObj.FirstName},
		{"last_name", updateObj.LastName},
	}}}
	_, err := m.Client.GetUserCollection().UpdateOne(ctx, query, update)
	if err != nil {
		return nil, errors.New("error in updating field")
	}
	user, err := m.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
