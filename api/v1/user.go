package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
)

func (server *Server) registerUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userToCreate := db.User{
		Email:    credentials.Email,
		Password: credentials.Password,
		IsActive: false,
	}

	_, err := server.store.CreateUser(c, userToCreate)
	if err != nil {
		NewErrorNotAcceptable(err).Send(ctx)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var credentials authUserBody
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		NewErrorBadRequest(err).Send(ctx)
		return
	}

	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := server.store.GetUserByEmail(c, credentials.Email)
	if err != nil {
		NewErrorNotFound(err).Send(ctx)
		return
	}

	err = util.CheckPassword(credentials.Password, user.PasswordHash)
	if err != nil {
		NewErrorUnauthorized(err).Send(ctx)
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.accessTokenDuration)
	if err != nil {
		NewErrorInternalSeverError(err).Send(ctx)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.refreshTokenDuration)
	if err != nil {
		NewErrorInternalSeverError(err).Send(ctx)
		return
	}

	sessionID, err := server.store.CreateSession(c, db.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIP:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiresAt,
	})

	if err != nil {
		NewErrorInternalSeverError(err).Send(ctx)
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
