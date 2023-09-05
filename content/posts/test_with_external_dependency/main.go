package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	ctx := context.Background()

	repo, err := NewMongoRepo(ctx, "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	user := &User{
		ID:       primitive.NewObjectID(),
		Name:     "John",
		Email:    "john@example.com",
		Password: "password",
	}

	err = repo.CreateUser(ctx, user)
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %+v\n", user)
}
