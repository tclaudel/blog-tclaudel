---
title: "Hexagonal architecture with Go"
date: 2023-09-20T21:41:18+02:00
draft: true
---

## Hexagonal architecture

Hexagonal architecture was invented by [Alistair Cockburn](https://alistair.cockburn.us/hexagonal-architecture/) 
and published in 2005.
Is a software architecture that aims to create loosely coupled application components with isolation between 
business logic and technical details. 

The hexagonal architecture is also known as the ports and adapters architecture,

### Why ?

If offers a number of advantages, such as :
- **Independence from framework**: your application is no longer directly dependent on external libraries.
- **Testability**: writing tests is greatly facilitated by the decoupling of dependencies.
- **Flexibility ans scalability**: the application is more flexible and scalable because it is not tied to a specific 
framework. It's easier to change the framework or add new features.

### Key Components

To achieve this, we will define 3 components:
- **Domain**: the business logic of the application, your entities, your business rules, etc.
- **Ports**: the interfaces that define how the domain interacts with the outside world.
- **Adapters**: the implementations of the ports.

## Domain

The domain is the core of the application. It contains the business logic, the entities, the business rules, etc.
It must not depend on anything. It must be completely independent of the outside world.  

**⚠️ NO INFRASTRUCTURE CODE IN THE DOMAIN. ⚠️**

By dependencies, We consider Web framework, Database clients, libraries etc.

![Domain](/images/content/hexagonal_architecture_with_go/domain.png)

He is the example of a model inside the core :

```bash
.
└── internal
    └── domain
        └── entity
            └── user.go
```


***./internal/domain/entity/user.go__***
```go
package entity

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/google/uuid"
)

var (
	// ErrInvalidUserID is returned when the user id is invalid
	ErrInvalidUserID = errors.New("invalid user id")
	// ErrInvalidEmail is returned when the email is invalid
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUsername is returned when the username is invalid
	ErrInvalidUsername = errors.New("invalid username")
)

type UserParams struct {
	ID       string
	Username string
	Email    string
}

// User represents a user
type User struct {
	id       uuid.UUID
	username string
	email    *mail.Address
}

// NewUser creates a new user
func NewUser(params UserParams) (*User, error) {
	userID, err := uuid.Parse(params.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUserID, err.Error())
	}

	emailAddress, err := mail.ParseAddress(params.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, err.Error())
	}

	if params.Username == "" {
		return nil, fmt.Errorf("%w: %s", ErrInvalidUsername, "username cannot be empty")
	}

	return &User{
		id:       userID,
		username: params.Username,
		email:    emailAddress,
	}, nil
}

func (u User) ID() uuid.UUID {
	return u.id
}

func (u User) Username() string {
	return u.username
}

func (u User) Email() *mail.Address {
	return u.email
}
```

## Ports

The ports defines the way the business logic interact with the outside world. They decouple domain and adapters wich 
hard the next component. In Golang interface are used for it. We define two kinds of ports : primary/secondary:

![Ports](/images/content/hexagonal_architecture_with_go/ports.png)

```bash
internal/ports
├── primary
│   └── user.go
└── secondary
    └── user.go
```


### Primary

They defines the uses for the application, it represents the uses cases and the scenarios  
They are implemented by the domain

***internal/ports/user.go***
```go
package primary

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
)

var (
	ErrCreatingUser      = errors.New("error creating user")
	ErrUserAlreadyExists = errors.New("error user already exists")
	ErrUpdatingUser      = errors.New("error updating user")
	ErrRemovingUser      = errors.New("error removing user")
	ErrRetrievingUser    = errors.New("error retrieving user")
	ErrRetrievingUsers   = errors.New("error retrieving users")
	ErrUserNotFound      = errors.New("error user not found")
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	RemoveUser(ctx context.Context, id uuid.UUID) error
	RetrieveUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	RetrieveUsers(ctx context.Context) ([]*entity.User, error)
}
```

example of an implementation in the domain :

```internal/
├── domain
│   ├── entity
│   │   └── user.go
│   └── services
│       ├── user_create.go
│       ├── user_create_test.go
│       ├── user_remove.go
│       ├── user_remove_test.go
│       ├── user_retrieve.go
│       ├── user_retrieve_test.go
│       ├── users_retrieve.go
│       ├── users_retrieve_test.go
│       ├── user_update.go
│       ├── user_update_test.go
│       └── user_use_case.go
└── ports
    ├── primary
    │   └── user.go
    └── secondary
        └── user.go
```

***domain/services/user_use_case.go***
```go
package services

import "github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"

type UserUseCase struct {
	userRepository secondary.UserRepository
}

func NewUserUseCase(userRepository secondary.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepository: userRepository,
	}
}
```

***domain/services/user_create.go***
```go
package services

import (
	"context"
	"errors"
	"fmt"

	"log"

	"github.com/tclaudel/hexagonal-architecture-golang/internal/domain/entity"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/primary"
	"github.com/tclaudel/hexagonal-architecture-golang/internal/ports/secondary"
)

func (u UserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	err := u.userRepository.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error saving user: %s", err.Error())

		if errors.Is(err, secondary.ErrUserAlreadyExists) {
			return fmt.Errorf("%w: %s", primary.ErrUserAlreadyExists, err.Error())
		}

		return fmt.Errorf("%w: %s", primary.ErrCreatingUser, err.Error())
	}

	return nil
}
```
