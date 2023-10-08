---
title: "Golang: Testing Unreachable Dependencies"
date: 2023-08-29T08:28:50+02:00
draft: false
---
## The Humble Object Design Pattern: Navigating the Challenges of Unit Testing


In the evolving landscape of software development, the need for code reliability and maintainability stands paramount. A
consistent hurdle developers encounter is the intricate task of unit testing code entangled with external components. To
address this, the software community has introduced a promising solution: The Humble Object design pattern.

### The Problem: Code Intertwined with External Dependencies

Consider the scenario of developing software where your logic is deeply embedded with user interfaces, databases,
filesystems, or other third-party elements. Such tight coupling of business logic with these external systems often
complicates unit testing due to:

1. **Non-deterministic Behavior**: External systems can exhibit unpredictable behavior, rendering the establishment of
   consistent unit tests challenging.
2. **Test Environment Complexity**: Configuring test environments that replicate real-world external systems can be
   intricate and prone to errors.
3. **Slow Tests**: Direct interactions with live systems, such as databases, can decelerate the testing process,
   causing disruptions in the development workflow.

### The Humble Object Solution

The core principle of the Humble Object pattern is the strategic decoupling of intricate logic from external
interactions. This bifurcation translates into:

1. **The Humble Object**: This component, kept intentionally simple, encapsulates interactions with external systems.
   Given its minimal business logic, it might not always be the focus of rigorous unit tests.
2. **The Logic Object**: This is where the core business logic resides. Stripped of external ties, it becomes inherently
   more testable, paving the way for comprehensive unit tests that validate the software's intended behavior.

### Benefits at a Glance

- **Focused Testing**: By distinguishing logic from external touchpoints, developers can zero in on testing the
  fundamental business logic without external distractions.
- **Enhanced Maintainability**: The demarcation between logic and external dependencies refines code clarity and eases
  maintenance efforts.
- **Scalability**: As the software expands, this clear separation ensures that integrating new external components or
  modifying existing ones doesn't call for extensive refactoring.

## The Humble Object in Action

Let's explore the Humble Object pattern by building a simple application that fetches data from a database and
performs some business logic on it. We'll then test this application using the Humble Object pattern.

### The Application

Our application will be a simple Go program that fetches data from a database and performs some business logic on it.

### The mongo client

We'll use a simple MongoDB database for our application. The database will have a single collection, `users`, with the
following struct:

***mongo.go***

```go
package main

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrConnectingToMongoDatabase = errors.New("error connecting to mongo database")
	ErrInsertingUser             = errors.New("error inserting user")
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email,omitempty"`
	Password string             `bson:"password,omitempty"`
}

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(ctx context.Context, mongoURI string) (*MongoRepo, error) {
	const (
		dbName         = "test"
		collectionName = "users"
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectingToMongoDatabase, err)
	}

	collection := client.Database(dbName).Collection(collectionName)

	return &MongoRepo{
		collection: collection,
	}, nil
}

func (m *MongoRepo) CreateUser(ctx context.Context, user *User) error {
	_, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInsertingUser, err)
	}

	return nil
}
```

### The Application

***main.go***
```go
package main

import (
	"context"

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
}
```

## Let's test

You've seen that our application works, but how do we test it? We can't just run the application and see if it works. 
We don't want to connect to a real database in our unit tests.

### The Problem

Let's start by writing a test for the `CreateUser` function:

***mongo_test.go***

```go
package main

import (
    "context"
    "testing"
)

func TestMongoRepo_CreateUser(t *testing.T) {
    ctx := context.Background()

    repo, err := NewMongoRepo(ctx, "mongodb://localhost:27017")
    if err != nil {
        t.Fatalf("error creating mongo repo: %s", err)
    }

    user := &User{
        ID:         primitive.NewObjectID(),
        Name:     "John",
        Email:    "john@example.com",
        Password: "password",
    }

    err = repo.CreateUser(ctx, user)
    if err != nil {
        t.Fatalf("error creating user: %s", err)
    }
}
```

If you run the test, you'll see that it fails:

```bash
$ go test -v -run TestMongoRepo_CreateUser
=== RUN   TestMongoRepo_CreateUser
    main_test.go:27: error creating user: error inserting user: server selection error: server selection timeout, current topology: { Type: Unknown, Servers: [{ Addr: localhost:27017, Type: Unknown, Average RTT: 0, Last error: connection() error occured during connection handshake: dial tcp 127.0.0.1:27017: connect: connection refused }, ] }
--- FAIL: TestMongoRepo_CreateUser (30.00s)
FAIL
exit status 1
FAIL    github.com/tclaudel/blog-tclaudel/content/posts/test_with_external_dependency   30.008s
```

The test fails because it can't connect to the database. We don't want to connect to a real database in our unit tests.
We need to mock the database.

### Decoupling With Interfaces

To mock the database, we need to decouple our code from the database calls. We can do this by using an interface, let's 
see which calls we need to mock:

```go
func (m *MongoRepo) CreateUser(ctx context.Context, user *User) error {
    _, err := m.collection.InsertOne(ctx, user) // <-- Database call
    if err != nil {
        return fmt.Errorf("%w: %s", ErrInsertingUser, err)
    }

    return nil
}
```

We can create an interface that contains these two methods:

***mongo.go***
```go
type MongoCaller interface {
    InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (
        *mongo.InsertOneResult, error)
}
```

We can then change our `MongoRepo` to use this struct instead of the `mongo.Collection`:
***mongo.go***
```go
type MongoRepo struct {
    mongoCaller MongoCaller
}
```
Modify the NewMongoRepo function to use the `mongoCaller` field:
***mongo.go***
```go
// [...]
func NewMongoRepo(ctx context.Context, mongoURI string) (*MongoRepo, error) {
    const (
        dbName         = "test"
        collectionName = "users"
    )
    
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        return nil, fmt.Errorf("%w: %s", ErrConnectingToMongoDatabase, err)
    }
    
    collection := client.Database(dbName).Collection(collectionName)
    
    return &MongoRepo{
        mongoCaller: collection,
    }, nil
}
// [...]
```
You have nothing to change in the `NewMongoRepo` function, it still returns a `MongoRepo` that implements the
`MongoCaller` interface.

Let's create a mock for this interface:

We implement the `InsertOne` method. Our Implementation will be very simple we will just return the ID of
the user we want to insert:

***mongo_mock.go***
```go
package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	return &mongo.InsertOneResult{
		InsertedID: doc.ID,
	}, nil
}
```

We can now use this mock in our test:

***mongo_test.go***
```go
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
```

If you run the test, you'll see that it passes:

```bash
$ go test -v -run TestMongoRepo_CreateUser
go test -v -run TestMongoRepo_CreateUser
=== RUN   TestMongoRepo_CreateUser
--- PASS: TestMongoRepo_CreateUser (0.00s)
PASS
ok      github.com/tclaudel/blog-tclaudel/content/posts/test_with_external_dependency   0.005s
```

Let's implement an errorous mock to make our test fail:

***mongo_mock.go***
```go
package main

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	emailWhichTriggersError = "error@error.com"
)

var _ MongoCaller = (*MockMongo)(nil)

type MockMongo struct{}

func NewMockMongo() *MongoRepo {
	return &MongoRepo{
		humbleObject: HumbleObject{
			mongoCaller: MockMongo{},
		},
	}
}

func (m MockMongo) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (
	*mongo.InsertOneResult, error,
) {
	doc, ok := document.(*User)
	if !ok {
		return nil, ErrInsertingUser
	}

	if doc.Email == emailWhichTriggersError {
		return nil, ErrInsertingUser
	}

	return &mongo.InsertOneResult{
		InsertedID: doc.ID,
	}, nil
}
```

Let's write another test to make sure that our code handles errors correctly:

***mongo_test.go***
```go
 // [...]
func TestMongoRepo_CreateUserError(t *testing.T) {
    ctx := context.Background()
    
    repo := NewMockMongo()
    
    user := &User{
        ID:       primitive.NewObjectID(),
        Name:     "John",
        Email:    emailWhichTriggersError,
        Password: "password",
    }
    
    err := repo.CreateUser(ctx, user)
    assert.ErrorIs(t, err, ErrInsertingUser)
}
```
Run the test :
```bash
$ go test -v -run TestMongoRepo_CreateUserError
=== RUN   TestMongoRepo_CreateUserError
--- PASS: TestMongoRepo_CreateUserError (0.00s)
PASS
ok      github.com/tclaudel/blog-tclaudel/content/posts/test_with_external_dependency   0.005s
```

## Conclusion

The Humble Object pattern is a powerful tool for decoupling intricate logic from external dependencies. This separation
allows developers to focus on testing the core business logic without the distractions of external systems. The
resulting tests are more reliable, maintainable, and scalable.