package db

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	PageID   int64 `form:"page_id" json:"page_id" binding:"min=1"`
	PageSize int64 `form:"page_id" json:"page_size" binding:"min=0,max=500"`
}

func (pagination *Pagination) getFindOptions() *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSkip((pagination.PageID - 1) * pagination.PageSize)
	findOptions.SetLimit(pagination.PageSize)

	return findOptions
}

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	ID           string `bson:"_id" json:"id"`
	Email        string `bson:"email" json:"email" binding:"required"`
	PasswordHash string `bson:"passwordHash" json:"passwordHash"`
	Password     string `bson:"password" json:"password" binding:"required"`
	Role         Role   `bson:"role" json:"role,omitempty"`
	CreatedAt    int64  `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt   int64  `bson:"modifiedAt" json:"modifiedAt,omitempty"`
}

type Author struct {
	ID           string `bson:"_id" json:"id"`
	FirstName    string `bson:"firstName" json:"firstName,omitempty"`
	LastName     string `bson:"lastName" json:"lastName,omitempty"`
	Name         string `bson:"name" json:"name" binding:"required"`
	WebsiteURL   string `bson:"websiteUrl" json:"websiteUrl,omitempty"`
	InstagramURL string `bson:"instagramUrl" json:"instagramUrl,omitempty"`
	YoutubeURL   string `bson:"youtubeUrl" json:"youtubeUrl,omitempty"`
	ImageBase64  string `bson:"imageBase64" json:"imageBase64,omitempty"`
	RecipeCount  int    `bson:"recipeCount" json:"recipeCount,omitempty"`
	UserID       string `bson:"userId" json:"userId,omitempty"`
	UserCreated  User   `bson:"userCreated" json:"userCreated"`
	CreatedAt    int64  `bson:"createdAt" json:"createdAt"`
	ModifiedAt   int64  `bson:"modifiedAt" json:"modifiedAt"`
}

type Category string

const (
	Breakfast Category = "breakfast"
	Main      Category = "main"
	Desert    Category = "desert"
	Smoothie  Category = "smoothie"
	Baby      Category = "baby"
	Drink     Category = "drink"
)

type AmountUnit string

const (
	Milliliters AmountUnit = "ml"
	Liters      AmountUnit = "l"
	Milligrams  AmountUnit = "mg"
	Grams       AmountUnit = "g"
	Tablespoon  AmountUnit = "tbs"
	Teaspoon    AmountUnit = "tsp"
	Piece       AmountUnit = "pc"
)

type Ingredient struct {
	Name   string `bson:"name" json:"name"`
	Amount int    `bson:"amount" json:"amount"`
	Unit   string `bson:"unit" json:"unit"`
}

type PrepStep struct {
	Rank        int    `bson:"rank" json:"rank"`
	Description string `bson:"description" json:"description"`
}

type Recipe struct {
	ID          string       `bson:"_id" json:"id"`
	Name        string       `bson:"name" json:"name"`
	RecipeURL   string       `bson:"recipeUrl" json:"recipeUrl,omitempty"`
	TimeM       int          `bson:"timeM" json:"timeM"`
	Category    Category     `bson:"category" json:"category"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
	PrepSteps   []PrepStep   `bson:"prepSteps" json:"prepSteps"`
	AuthorID    string       `bson:"authorId" json:"authorId,omitempty" binding:"required"`
	Author      Author       `bson:"author" json:"author"`
	UserID      string       `bson:"userId" json:"userId,omitempty"`
	User        User         `bson:"user" json:"user"`
	CreatedAt   int64        `bson:"createdAt" json:"createdAt"`
	ModifiedAt  int64        `bson:"modifiedAt" json:"modifiedAt"`
}
