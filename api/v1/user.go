package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type authUserBody struct {
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
} // @name authUserBody

type registerResponse struct {
	SessionID             uuid.UUID `json:"sessionId"`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	UserID                string    `json:"userId"`
} // @name registerResponse

func (server *Server) registerUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err) // TODO: Create proper error response
		return
	}

	userColl := db.NewUserCollection(server.dbClient, server.dbName)

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userToCreate := db.User{
		Email:    credentials.Email,
		Password: credentials.Email,
	}

	_, err := userColl.CreateUser(c, userToCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err) // TODO: Create proper error response
		return
	}

	ctx.Status(http.StatusCreated)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err) // TODO: Create proper error response
		return
	}

	userColl := db.NewUserCollection(server.dbClient, server.dbName)

	// user, err := userColl.GetUserByEmail()
}
