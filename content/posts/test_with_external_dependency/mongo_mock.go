package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	emailWitchTriggersError = "error@error.com"
)

var _ MongoCaller = (*MockMongo)(nil)

type MockMongo struct{}

func NewMockMongo() *MongoRepo {
	return &MongoRepo{
		mongoCaller: MockMongo{},
	}
}

func (m MockMongo) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (
	*mongo.InsertOneResult, error,
) {
	doc, ok := document.(*User)
	if !ok {
		return nil, ErrInsertingUser
	}

	if doc.Email == emailWitchTriggersError {
		return nil, ErrInsertingUser
	}

	return &mongo.InsertOneResult{
		InsertedID: doc.ID,
	}, nil
}
