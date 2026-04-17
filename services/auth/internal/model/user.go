package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	Name           string             `bson:"name"`
	HashedPassword string             `bson:"hashed_password"`
	CreatedAt      int64              `bson:"created_at"`
}

// constructor now
func NewStore(collection *mongo.Collection) *Store {
	return &Store{
		collection: collection,
	}
}

// create user now
func (s *Store) CreateUser(ctx context.Context, email, name, hashedPassword string) (*User, error) {
	user := User{
		Email:          email,
		Name:           name,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now().Unix(),
	}

	//insert into the database
	res, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return &user, nil
}

func (s *Store) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	//check if the user exists
	var user User

	filter := bson.M{"email": email}
	err := s.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil

}


