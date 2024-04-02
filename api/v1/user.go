package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
)

type authUserBody struct {
	Email    string `json:"email,omitempty" binding:"required"`
	Password string `json:"password,omitempty" binding:"required,min=6"`
} // @name authUserBody

func (server *Server) registerUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err) // TODO: Create proper error response
		return
	}

	userColl := db.NewUserCollection(server.config.dbClient, server.config.dbName)

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userToCreate := db.User{
		Email:    credentials.Email,
		Password: credentials.Email,
		IsActive: false,
	}

	_, err := userColl.CreateUser(c, userToCreate)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err) // TODO: Create proper error response
		return
	}

	ctx.Status(http.StatusCreated)
}

type loginResponse struct {
	SessionID             string  `json:"sessionId"`
	AccessToken           string  `json:"accessToken"`
	AccessTokenExpiresAt  int64   `json:"accessTokenExpiresAt"`
	RefreshToken          string  `json:"refreshToken"`
	RefreshTokenExpiresAt int64   `json:"refreshTokenExpiresAt"`
	User                  db.User `json:"user"`
} // @name loginResponse

func (server *Server) loginUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, err) // TODO: Create proper error response
		return
	}

	userColl := db.NewUserCollection(server.config.dbClient, server.config.dbName)

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := userColl.GetUserByEmail(c, credentials.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, err) // TODO: Create proper error response
		return
	}

	err = util.CheckPassword(credentials.Password, user.PasswordHash)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, err)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.accessTokenDuration)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.refreshTokenDuration)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	sessionColl := db.NewSessionCollection(server.config.dbClient, server.config.dbName)

	sessionID, err := sessionColl.CreateSession(c, db.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIP:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	res := loginResponse{
		SessionID:             sessionID.Hex(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		User:                  user,
	}

	ctx.JSON(http.StatusAccepted, res)
}
