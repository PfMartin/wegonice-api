package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type registerUserBody struct {
	email    string `json:"email,omitempty"`
	password string `json:"password,omitempty"`
} // @name registerUserBody

type registerResponse struct {
	SessionID             uuid.UUID `json:""sessionId`
	AccessToken           string    `json:"accessToken"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt"`
	RefreshToken          string    `json:"refreshToken"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt"`
	UserID                string    `json:"userId"`
} // @name registerResponse

func (server *Server) registerUser(ctx *gin.Context) {
	var credentials registerUserBody
	if err := ctx.ShouldBindJSON(credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, nil) // TODO: Create proper error response
		return
	}

	userColl := db.NewUserCollection(server.dbClient, server.dbName)

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userToCreate := db.User{
		Email:    credentials.email,
		Password: credentials.email,
	}

	userID, err := userColl.CreateUser(c, userToCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil) // TODO: Create proper error response
		return
	}

	userIDString := userID.Hex()

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(userIDString, 10*time.Minute)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, nil) // TODO: Create proper error response
		return
	}

	res := registerResponse{
		SessionID:            accessPayload.ID,
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiresAt,
		// RefreshToken:
		// RefreshTokenExpiresAt:
		UserID: userIDString,
	}

	ctx.JSON(http.StatusOK, res)
}
