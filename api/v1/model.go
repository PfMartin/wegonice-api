package api

import "github.com/PfMartin/wegonice-api/db"

type authUserBody struct {
	Email    string `json:"email,omitempty" binding:"required"` //TODO: Email validation
	Password string `json:"password,omitempty" binding:"required,min=6"`
} // @name authUserBody

type loginResponse struct {
	SessionID             string  `json:"sessionId"`
	AccessToken           string  `json:"accessToken"`
	AccessTokenExpiresAt  int64   `json:"accessTokenExpiresAt"`
	RefreshToken          string  `json:"refreshToken"`
	RefreshTokenExpiresAt int64   `json:"refreshTokenExpiresAt"`
	User                  db.User `json:"user"`
} // @name loginResponse
