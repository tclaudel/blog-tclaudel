---
title: "Golang: Tester les Dépendances Inaccessibles"
date: 2023-08-29T08:28:50+02:00
draft: false
---
## Le Modèle de Conception Humble object : Naviguer à travers les Défis des Tests Unitaires

Dans le paysage évolutif du développement logiciel, le besoin de fiabilité et de maintenabilité du code est primordial.
Un obstacle constant auquel les développeurs sont confrontés est la tâche complexe des tests unitaires pour du code
imbriqué avec des composants externes. Pour remédier à cela, la communauté logicielle a introduit une solution
prometteuse : le modèle de conception *humble object*.

> Code d'exemple ici: https://github.com/tclaudel/blog-tclaudel/tree/main/content/posts/test_with_external_dependency

### Le Problème : Code Entremêlé avec des Dépendances Externes

Considérez le scénario du développement logiciel où votre logique est profondément intégrée avec des interfaces
utilisateur, des bases de données, des systèmes de fichiers ou d'autres éléments tiers. Un couplage étroit de la
logique métier avec ces systèmes externes complique souvent les tests unitaires en raison de :

1. **Comportement Non Déterministe**: Les systèmes externes peuvent présenter un comportement imprévisible, rendant
   l'établissement de tests unitaires cohérents difficile.
2. **Complexité de l'Environnement de Test**: Configurer des environnements de test qui reproduisent des systèmes
   externes du monde réel peut être complexe et sujet aux erreurs.
3. **Tests Lents**: Les interactions directes avec des systèmes en direct, tels que les bases de données, peuvent
   ralentir le processus de test, perturbant le flux de développement.

### La Solution du Modèle de Conception Humble Object

Le principe fondamental du modèle de conception humble object est le découplage stratégique de la logique complexe des
interactions externes. Cette bifurcation se traduit par :

1. **L'Objet Modeste (Humble object)**: Ce composant, maintenu intentionnellement simple, encapsule les interactions avec les systèmes
   externes. Étant donné sa logique métier minimale, il ne sera pas toujours au centre de tests unitaires rigoureux.
2. **L'Objet Logique** : C'est là que réside la logique métier principale. Dépourvu de liens externes, il devient
   intrinsèquement plus testable, ouvrant la voie à des tests unitaires complets qui valident le comportement voulu du
   logiciel.

### Avantages

1. **Tests Ciblés** : En distinguant la logique des points de contact externes, les développeurs peuvent se concentrer
   sur le test de la logique métier fondamentale sans distractions externes.
2. **Maintenabilité Améliorée** : La démarcation entre la logique et les dépendances externes affine la clarté du code
   et facilite les efforts de maintenance.
3. **Scalabilité** : À mesure que le logiciel se développe, cette séparation claire garantit que l'intégration de
   nouveaux composants externes ou la modification de ceux existants ne nécessite pas de refactoring étendu.

## Le Modèle de Conception Humble Object en Action

Explorons le modèle de conception *humble object* en construisant une application simple qui récupère des données d'une
base de données et effectue une certaine logique métier dessus. Nous testerons ensuite cette application en utilisant le
modèle de conception *humble object*.

### L'Application

Notre application sera un simple programme Go qui récupère des données d'une base de données et effectue une certaine
logique métier dessus.

### Le client MongoDB

Nous utiliserons une base de données MongoDB simple pour notre application. La base de données aura une seule
collection, `users`, avec la structure suivante :

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

### L'Application

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

## Passons aux Tests

Vous avez vu que notre application fonctionne, mais comment la tester ? Nous ne pouvons pas simplement exécuter
l'application et voir si elle fonctionne. Nous ne voulons pas nous connecter à une vraie base de données dans nos tests
unitaires.

### Le Problème

Commençons par écrire un test pour la fonction `CreateUser` :

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
		ID:       primitive.NewObjectID(),
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

Si vous exécutez le test, vous verrez qu'il échoue :

```bash
$ go test -v -run TestMongoRepo_CreateUser
=== RUN   TestMongoRepo_CreateUser
    main_test.go:27: error creating user: error inserting user: server selection error: server selection timeout, current topology: { Type: Unknown, Servers: [{ Addr: localhost:27017, Type: Unknown, Average RTT: 0, Last error: connection() error occured during connection handshake: dial tcp 127.0.0.1:27017: connect: connection refused }, ] }
--- FAIL: TestMongoRepo_CreateUser (30.00s)
FAIL
exit status 1
FAIL    github.com/tclaudel/blog-tclaudel/content/posts/test_with_external_dependency   30.008s
```

Le test échoue car il ne peut pas se connecter à la base de données. Nous ne voulons pas nous connecter à une vraie base
de données dans nos tests unitaires. Nous devons simuler la base de données.

### Découplage Avec des Interfaces

Pour simuler la base de données, nous devons découpler notre code des appels à la base de données. Nous pouvons le faire
en utilisant une interface, voyons quels appels nous devons simuler :

```go
func (m *MongoRepo) CreateUser(ctx context.Context, user *User) error {
_, err := m.collection.InsertOne(ctx, user) // <-- Database call
if err != nil {
return fmt.Errorf("%w: %s", ErrInsertingUser, err)
}

return nil
}
```

Nous pouvons créer une interface qui contient ces deux méthodes :

***mongo.go***

```go
type MongoCaller interface {
InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (
*mongo.InsertOneResult, error)
}
```

Nous pouvons ensuite modifier notre `MongoRepo` pour utiliser cette structure au lieu de `mongo.Collection` :

***mongo.go***

```go
type MongoRepo struct {
mongoCaller MongoCaller
}
```

Modifiez la fonction `NewMongoRepo` pour utiliser le champ `mongoCaller` :
***mongo.go***

```go
// [...]
func NewMongoRepo(ctx context.Context, mongoURI string) (*MongoRepo, error) {
const (
dbName = "test"
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

Vous n'avez rien à changer dans la fonction `NewMongoRepo`, elle retourne toujours un `MongoRepo` qui implémente
l'interface `MongoCaller`.

Créons maintenant un faux pour cette interface :

Nous implémentons la méthode `InsertOne`. Notre implémentation sera très simple, nous allons simplement retourner l'ID
de l'utilisateur que nous voulons insérer :

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

Nous pouvons maintenant utiliser ce faux dans notre test :

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

Si vous exécutez le test, vous verrez qu'il réussit :

```bash
$ go test -v -run TestMongoRepo_CreateUser
go test -v -run TestMongoRepo_CreateUser
=== RUN   TestMongoRepo_CreateUser
--- PASS: TestMongoRepo_CreateUser (0.00s)
PASS
ok      github.com/tclaudel/blog-tclaudel/content/posts/test_with_external_dependency   0.005s
```

Implémentons maintenant un faux erroné pour que notre test échoue :

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

Écrivons un autre test pour nous assurer que notre code gère correctement les erreurs :

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

Le modèle de conception *humble object* est un outil puissant pour découpler la logique complexe des dépendances externes.
Cette séparation permet aux développeurs de se concentrer sur le test de la logique

> Code d'exemple ici: https://github.com/tclaudel/blog-tclaudel/tree/main/content/posts/test_with_external_dependency

### Sources :
- https://martinfowler.com/bliki/HumbleObject.html
- http://xunitpatterns.com/Humble%20Object.html
