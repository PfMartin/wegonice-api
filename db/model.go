package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	PageID   int64 `form:"page_id" json:"page_id" binding:"required,min=1"`
	PageSize int64 `form:"page_size" json:"page_size" binding:"required,min=0,max=500"`
}

func (pagination *Pagination) getFindOptions() *options.FindOptions {
	findOptions := options.Find()
	findOptions.SetSkip((pagination.PageID - 1) * pagination.PageSize)
	findOptions.SetLimit(pagination.PageSize)

	return findOptions
}

func (pagination *Pagination) getSkipStage() bson.M {
	return bson.M{"$skip": (pagination.PageID - 1) * pagination.PageSize}
}

func (pagination *Pagination) getLimitStage() bson.M {
	return bson.M{"$limit": pagination.PageSize}
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
	IsActive     bool   `bson:"isActive" json:"isActive,omitempty"`
	CreatedAt    int64  `bson:"createdAt" json:"createdAt,omitempty"`
	ModifiedAt   int64  `bson:"modifiedAt" json:"modifiedAt,omitempty"`
}

type Session struct {
	ID           string `bson:"_id" json:"id"`
	UserID       string `bson:"userId" json:"userId"`
	User         User   `bson:"user" json:"user"`
	RefreshToken string `bson:"refreshToken" json:"refreshToken"`
	UserAgent    string `bson:"userAgent" json:"userAgent"`
	IsBlocked    bool   `bson:"isBlocked" json:"isBlocked"`
	ClientIP     string `bson:"clientIp" json:"clientIp"`
	ExpiresAt    int64  `bson:"expiresAt" json:"expiresAt"`
}

type Author struct {
	ID           string `bson:"_id" json:"id"`
	FirstName    string `bson:"firstName" json:"firstName,omitempty"`
	LastName     string `bson:"lastName" json:"lastName,omitempty"`
	Name         string `bson:"name" json:"name" binding:"required"`
	WebsiteURL   string `bson:"websiteUrl" json:"websiteUrl,omitempty"`
	InstagramURL string `bson:"instagramUrl" json:"instagramUrl,omitempty"`
	YoutubeURL   string `bson:"youtubeUrl" json:"youtubeUrl,omitempty"`
	ImageName    string `bson:"imageName" json:"imageName,omitempty"`
	RecipeCount  int    `bson:"recipeCount" json:"recipeCount,omitempty"`
	UserID       string `bson:"userId" json:"userId,omitempty"`
	UserCreated  User   `bson:"userCreated" json:"userCreated"`
	CreatedAt    int64  `bson:"createdAt" json:"createdAt"`
	ModifiedAt   int64  `bson:"modifiedAt" json:"modifiedAt"`
} // @name Author

type AuthorToCreate struct {
	FirstName    string `bson:"firstName" json:"firstName,omitempty" example:"Moe"`
	LastName     string `bson:"lastName" json:"lastName,omitempty" example:"Zarella"`
	Name         string `bson:"name" json:"name" binding:"required" example:"Moe Zarella"`
	WebsiteURL   string `bson:"websiteUrl" json:"websiteUrl,omitempty" example:"https://www.moezarella.com"`
	InstagramURL string `bson:"instagramUrl" json:"instagramUrl,omitempty" example:"https://wwww.instagram.com/moezarella/"`
	YoutubeURL   string `bson:"youtubeUrl" json:"youtubeUrl,omitempty" example:"https://www.youtube.com/channel/UCy8asdgasdf7RcC6OZffZA"`
	ImageName    string `bson:"imageName" json:"imageName,omitempty" example:"moezarella.png"`
	UserID       string `bson:"userId" json:"userId" binding:"required" example:"660c4b99bc1bc4aabe3e6cd1"`
} // @name AuthorToCreate

type AuthorUpdate struct {
	FirstName    string `bson:"firstName" json:"firstName,omitempty" example:"Moe"`
	LastName     string `bson:"lastName" json:"lastName,omitempty" example:"Zarella"`
	Name         string `bson:"name" json:"name" example:"Moe Zarella"`
	WebsiteURL   string `bson:"websiteUrl" json:"websiteUrl,omitempty" example:"https://www.moezarella.com"`
	InstagramURL string `bson:"instagramUrl" json:"instagramUrl,omitempty" example:"https://wwww.instagram.com/moezarella/"`
	YoutubeURL   string `bson:"youtubeUrl" json:"youtubeUrl,omitempty" example:"https://www.youtube.com/channel/UCy8asdgasdf7RcC6OZffZA"`
	ImageName    string `bson:"imageName" json:"imageName,omitempty" example:"moezarella.png"`
} // @name AuthorUpdate

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
	Name   string     `bson:"name" json:"name" example:"flour"`
	Amount int        `bson:"amount" json:"amount" example:"100"`
	Unit   AmountUnit `bson:"unit" json:"unit" example:"g"`
}

type PrepStep struct {
	Rank        int    `bson:"rank" json:"rank" example:"1"`
	Description string `bson:"description" json:"description" example:"Dice the onions"`
}

type Recipe struct {
	ID          string       `bson:"_id" json:"id"`
	Name        string       `bson:"name" json:"name"`
	ImageName   string       `bson:"imageName" json:"imageName,omitempty"`
	RecipeURL   string       `bson:"recipeUrl" json:"recipeUrl,omitempty"`
	TimeM       int          `bson:"timeM" json:"timeM"`
	Category    Category     `bson:"category" json:"category"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
	PrepSteps   []PrepStep   `bson:"prepSteps" json:"prepSteps"`
	AuthorID    string       `bson:"authorId" json:"authorId,omitempty" binding:"required"`
	Author      Author       `bson:"author" json:"author"`
	UserID      string       `bson:"userId" json:"userId,omitempty"`
	UserCreated User         `bson:"userCreated" json:"userCreated"`
	CreatedAt   int64        `bson:"createdAt" json:"createdAt"`
	ModifiedAt  int64        `bson:"modifiedAt" json:"modifiedAt"`
}

type RecipeToCreate struct {
	Name        string       `bson:"name" json:"name" example:"Pancakes" binding:"required"`
	ImageName   string       `bson:"imageName" json:"imageName,omitempty" example:"Pancakes.png"`
	RecipeURL   string       `bson:"recipeUrl" json:"recipeUrl,omitempty" example:"https://www.allthepancakes.com/pancakes"`
	TimeM       int          `bson:"timeM" json:"timeM" example:"30"`
	Category    Category     `bson:"category" json:"category" example:"breakfast"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients"`
	PrepSteps   []PrepStep   `bson:"prepSteps" json:"prepSteps"`
	AuthorID    string       `bson:"authorId" json:"authorId,omitempty" binding:"required" example:"660c4b99bc1bc4aabe126cd1"`
	UserID      string       `bson:"userId" json:"userId,omitempty" binding:"required" example:"660c4b99bc1bc4aabe126cd1"`
} // @name RecipeToCreate

type RecipeUpdate struct {
	Name        string       `bson:"name" json:"name,omitempty" example:"Pancakes"`
	ImageName   string       `bson:"imageName" json:"imageName,omitempty" example:"Pancakes.png"`
	RecipeURL   string       `bson:"recipeUrl" json:"recipeUrl,omitempty" example:"https://www.allthepancakes.com/pancakes"`
	TimeM       int          `bson:"timeM" json:"timeM,omitempty" example:"30"`
	Category    Category     `bson:"category" json:"category,omitempty" example:"breakfast"`
	Ingredients []Ingredient `bson:"ingredients" json:"ingredients,omitempty"`
	PrepSteps   []PrepStep   `bson:"prepSteps" json:"prepSteps,omitempty"`
	AuthorID    string       `bson:"authorId" json:"authorId,omitempty" example:"660c4b99bc1bc4aabe126cd1"`
} // @name RecipeUpdate
