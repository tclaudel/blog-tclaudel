package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoRepo_CreateUser(t *testing.T) {
	ctx := context.Background()

	repo := NewMockMongo()

	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("error creating user: %s", err)
	}
}

func TestMongoRepo_CreateUserError(t *testing.T) {
	ctx := context.Background()

	repo := NewMockMongo()

	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "John",
		Email:    emailWitchTriggersError,
		Password: "password",
	}

	err := repo.CreateUser(ctx, user)
	assert.ErrorIs(t, err, ErrInsertingUser)
}
