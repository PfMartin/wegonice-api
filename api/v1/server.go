package api

import (
	"time"

	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     ServerConfig
	store      db.DBStore
	tokenMaker token.Maker
	router     *gin.Engine
}

type ServerConfig struct {
	url                  string
	basePath             string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
	corsAllowedOrigins   []string
}

func NewServer(
	store db.DBStore,
	url string,
	basePath string,
	tokenSymmetricKey string,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
	corsAllowedOrigins []string,
) *Server {
	tokenMaker, err := token.NewPasetoMaker(tokenSymmetricKey)
	if err != nil {
		log.Err(err).Msg("cannot create token maker")
		return nil
	}

	config := ServerConfig{
		url:                  url,
		basePath:             basePath,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		corsAllowedOrigins:   corsAllowedOrigins,
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRoutes()

	return server
}

func (server *Server) setupRoutes() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: server.config.corsAllowedOrigins,
		AllowMethods: []string{"GET", "POST", "PATCH", "DELETE"},
	}))

	v1Routes := router.Group(server.config.basePath)

	v1Routes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1Routes.GET("/heartbeat", server.getHeartbeat)

	authRoutes := v1Routes.Group("/auth")
	authRoutes.POST("/register", server.registerUser)
	authRoutes.POST("/login", server.loginUser)

	authorRoutes := v1Routes.Group("/authors")
	authorRoutes.Use(authMiddleware(server.tokenMaker))
	authorRoutes.GET("", server.listAuthors)
	authorRoutes.POST("/", server.createAuthor)
	authorRoutes.GET("/:id", server.getAuthorByID)
	authorRoutes.PATCH("/:id", server.patchAuthorByID)
	authorRoutes.DELETE("/:id", server.deleteAuthorByID)

	server.router = router
}

func (server *Server) Start() error {
	return server.router.Run(server.config.url)
}
