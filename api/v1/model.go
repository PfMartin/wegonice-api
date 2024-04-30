package api

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
