package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DBStore interface {
	CreateUser(ctx context.Context, user User) (primitive.ObjectID, error)
	GetAllUsers(ctx context.Context, pagination Pagination) ([]User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	UpdateUserByID(ctx context.Context, userID string, userUpdate User) (int64, error)
	DeleteUserByID(ctx context.Context, userID string) (int64, error)

	CreateAuthor(ctx context.Context, author AuthorToCreate) (primitive.ObjectID, error)
	GetAllAuthors(ctx context.Context, pagination Pagination) ([]Author, error)
	GetAuthorByID(ctx context.Context, authorID string) (Author, error)
	UpdateAuthorByID(ctx context.Context, authorID string, authorUpdate AuthorUpdate) (int64, error)
	DeleteAuthorByID(ctx context.Context, authorID string) (int64, error)

	CreateRecipe(ctx context.Context, recipe RecipeToCreate) (primitive.ObjectID, error)
	GetAllRecipes(ctx context.Context, pagination Pagination) ([]Recipe, error)
	GetRecipeByID(ctx context.Context, recipeID string) (Recipe, error)
	UpdateRecipeByID(ctx context.Context, recipeID string, recipeUpdate Recipe) (int64, error)
	DeleteRecipeByID(ctx context.Context, recipeID string) (int64, error)

	CreateSession(ctx context.Context, session Session) (primitive.ObjectID, error)
	GetSessionByID(ctx context.Context, sessionID string) (Session, error)
}

type MongoDBStore struct {
	userCollection    *mongo.Collection
	authorCollection  *mongo.Collection
	recipeCollection  *mongo.Collection
	sessionCollection *mongo.Collection
}

func NewMongoDBStore(dbName, dbUser, dbPassword, dbURI string) *MongoDBStore {
	client, _ := NewDatabaseClient(dbName, dbUser, dbPassword, dbURI)

	database := client.Database(dbName)

	return &MongoDBStore{
		userCollection:    database.Collection("users"),
		authorCollection:  database.Collection("authors"),
		recipeCollection:  database.Collection("recipes"),
		sessionCollection: database.Collection("sessions"),
	}
}
