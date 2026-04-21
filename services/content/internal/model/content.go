package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Content struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	ContentType string             `bson:"content_type"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	ImageURL    string             `bson:"image_url"`
	Content     string             `bson:"content"`
	CreatedAt   int64              `bson:"created_at"`
}

// store struct
type Store struct {
	collection *mongo.Collection
}

// constructor
func NewStore(collection *mongo.Collection) *Store {
	return &Store{
		collection: collection,
	}
}

// AddContent to add the content to db
func (s *Store) AddContent(ctx context.Context, content *Content) (*Content, error) {
	//since we are already getting req in Content struct format only
	//go ahead and insert into the db
	res, err := s.collection.InsertOne(ctx, content)
	if err != nil {
		return nil, err
	}
	content.ID = res.InsertedID.(primitive.ObjectID)
	return content, nil
}

// GetContentByID to get the content by its id
func (s *Store) GetContentByUserId(ctx context.Context, userId string) ([]*Content, error) {
	//make a slice/array of contents
	var contents []*Content
	//make a filter
	filter := bson.M{"user_id": userId}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &contents); err != nil {
		return nil, err
	}
	return contents, nil
}

// DeleteContentById to delete content by its id
func (s *Store) DeleteContentByID(ctx context.Context, id string, userId string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	//make a filter
	filter := bson.M{"_id": objectId, "user_id": userId}
	_, err = s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}
