package api

import (
	"fmt"
	"strings"

	"github.com/PfMartin/wegonice-api/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			NewErrorUnauthorized(fmt.Errorf("authorization header is not provided")).Send(ctx)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			NewErrorUnauthorized(fmt.Errorf("invalid authorization header format")).Send(ctx)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			NewErrorUnauthorized(fmt.Errorf("unsupported authorization type %s", authorizationType)).Send(ctx)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			NewErrorUnauthorized(err).Send(ctx)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
