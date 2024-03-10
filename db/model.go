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
