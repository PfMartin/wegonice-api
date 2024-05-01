package api

import (
	"context"
	"net/http"
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
)

// registerUser
//
// @Summary 		Registers a user
// @Description A new user is registered by setting an unique email address and a password. The admins have to approve and active the user manually.
// @ID					auth-register-user
// @Tags				auth
// @Accept			json
// @Produce			json
// @Param				userCredentials 	body				authUserBody				true		"Credentials for the registration"
// @Success			201
// @Failure			400								{object}		ErrorBadRequest							"Bad Request"
// @Failure 		406								{object}		ErrorNotAcceptable					"Not Acceptable"
// @Router			/auth/register		[post]
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

// loginUser
//
// @Summary 		Logs a user in
// @Description A registered user is logged in with their email and matching password.
// @ID					auth-login-user
// @Tags				auth
// @Accept			json
// @Produce			json
// @Param				userCredentials 	body				authUserBody							true		"Credentials for the login"
// @Success			200								{object}		loginResponse											"Login response with required tokens"
// @Failure			400								{object}		ErrorBadRequest										"Bad Request"
// @Failure			401								{object}		ErrorUnauthorized									"Unauthorized"
// @Failure			404								{object}		ErrorNotFound											"Not Found"
// @Failure 		406								{object}		ErrorNotAcceptable								"Not Acceptable"
// @Failure 		500								{object}		ErrorInternalServerError					"Internal Server Error"
// @Router			/auth/login				[post]
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
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Email, server.config.refreshTokenDuration)
	if err != nil {
		NewErrorInternalServerError(err).Send(ctx)
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
		NewErrorInternalServerError(err).Send(ctx)
		return
	}

	res := loginResponse{
		SessionID:             sessionID.Hex(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiresAt,
		UserEmail:             user.ID,
	}

	ctx.JSON(http.StatusAccepted, res)
}
