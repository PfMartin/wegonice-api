package api

import "github.com/PfMartin/wegonice-api/db"

type getByIDRequest struct {
	ID string `uri:"id" binding:"required"`
}

type authUserBody struct {
	Email    string `json:"email,omitempty" binding:"required" example:"user@example.com"` //TODO: Email validation
	Password string `json:"password,omitempty" binding:"required,min=6" example:"s3cr3tP@ssw0rd"`
} // @name authUserBody

type loginResponse struct {
	SessionID             string `json:"sessionId" example:"660c4b99bc1bc4aabe3e6cd1"`
	AccessToken           string `json:"accessToken" example:"v2.local.example-session-token"`
	AccessTokenExpiresAt  int64  `json:"accessTokenExpiresAt" example:"1714462120"`
	RefreshToken          string `json:"refreshToken" example:"v2.local.example-refresh-token"`
	RefreshTokenExpiresAt int64  `json:"refreshTokenExpiresAt" example:"1714462120"`
	UserEmail             string `json:"userEmail" example:"user@example.com"`
} // @name loginResponse

type UserResponse struct {
	ID         string  `bson:"_id" json:"id" example:"660c4b99bc1bc4aabe3e6cd1"`
	Email      string  `bson:"email" json:"email" binding:"required" example:"user@example.com"`
	Role       db.Role `bson:"role" json:"role,omitempty" example:"user"`
	IsActive   bool    `bson:"isActive" json:"isActive,omitempty" example:"true"`
	CreatedAt  int64   `bson:"createdAt" json:"createdAt,omitempty" example:"1714462120"`
	ModifiedAt int64   `bson:"modifiedAt" json:"modifiedAt,omitempty" example:"1714462120"`
} // @name UserResponse

type AuthorBody struct {
	FirstName    string `bson:"firstName" json:"firstName,omitempty"`
	LastName     string `bson:"lastName" json:"lastName,omitempty"`
	Name         string `bson:"name" json:"name" binding:"required"`
	WebsiteURL   string `bson:"websiteUrl" json:"websiteUrl,omitempty"`
	InstagramURL string `bson:"instagramUrl" json:"instagramUrl,omitempty"`
	YoutubeURL   string `bson:"youtubeUrl" json:"youtubeUrl,omitempty"`
	ImageName    string `bson:"imageName" json:"imageName,omitempty"`
	UserID       string `bson:"userId" json:"userId,omitempty" binding:"required"`
} // @name AuthorBody

type AuthorResponse struct {
	ID           string       `bson:"_id" json:"id" example:"660c4b99bc1bc4aabe126cd1"`
	FirstName    string       `bson:"firstName" json:"firstName,omitempty" example:"Moe"`
	LastName     string       `bson:"lastName" json:"lastName,omitempty" example:"Zarella"`
	Name         string       `bson:"name" json:"name" binding:"required" example:"Moe Zarella"`
	WebsiteURL   string       `bson:"websiteUrl" json:"websiteUrl,omitempty" example:"https://www.moezarella.com"`
	InstagramURL string       `bson:"instagramUrl" json:"instagramUrl,omitempty" example:"https://wwww.instagram.com/moezarella/"`
	YoutubeURL   string       `bson:"youtubeUrl" json:"youtubeUrl,omitempty" example:"https://www.youtube.com/channel/UCy8asdgasdf7RcC6OZffZA"`
	ImageName    string       `bson:"imageName" json:"imageName,omitempty" example:"moezarella.png"`
	RecipeCount  int          `bson:"recipeCount" json:"recipeCount,omitempty" example:"5"`
	UserID       string       `bson:"userId" json:"userId,omitempty" example:"660c4b99bc1bc4aabe3e6cd1"`
	UserCreated  UserResponse `bson:"userCreated" json:"userCreated"`
	CreatedAt    int64        `bson:"createdAt" json:"createdAt" example:"1714462120"`
	ModifiedAt   int64        `bson:"modifiedAt" json:"modifiedAt" example:"1714462120"`
} // @name AuthorResponse

type RecipeResponse struct {
	ID          string          `bson:"_id" json:"id" example:"660c4b99bc1bc4aabe126cd1"`
	Name        string          `bson:"name" json:"name" example:"Pancakes"`
	ImageName   string          `bson:"imageName" json:"imageName,omitempty" example:"Pancakes.png"`
	RecipeURL   string          `bson:"recipeUrl" json:"recipeUrl,omitempty" example:"https://www.allthepancakes.com/pancakes"`
	TimeM       int             `bson:"timeM" json:"timeM" example:"30"`
	Category    db.Category     `bson:"category" json:"category" example:"breakfast"`
	Ingredients []db.Ingredient `bson:"ingredients" json:"ingredients"`
	PrepSteps   []db.PrepStep   `bson:"prepSteps" json:"prepSteps"`
	AuthorID    string          `bson:"authorId" json:"authorId,omitempty" binding:"required" example:"660c4b99bc1bc4aabe126cd1"`
	Author      AuthorResponse  `bson:"author" json:"author"`
	UserID      string          `bson:"userId" json:"userId,omitempty" example:"660c4b99bc1bc4aabe126cd1"`
	UserCreated UserResponse    `bson:"userCreated" json:"userCreated"`
	CreatedAt   int64           `bson:"createdAt" json:"createdAt" example:"1714462120"`
	ModifiedAt  int64           `bson:"modifiedAt" json:"modifiedAt" example:"1714462120"`
} // @name RecipeResponse
