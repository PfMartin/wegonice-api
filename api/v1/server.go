package api

import (
	"github.com/PfMartin/wegonice-api/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	dbClient   *mongo.Client
	dbName     string
	url        string
	basePath   string
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(dbClient *mongo.Client, dbName string, url string, basePath string, tokenSymmetricKey string) *Server {
	tokenMaker, err := token.NewPasetoMaker(tokenSymmetricKey)
	if err != nil {
		log.Err(err).Msg("cannot create token maker")
		return nil
	}

	server := &Server{
		dbClient:   dbClient,
		dbName:     dbName,
		url:        url,
		basePath:   basePath,
		tokenMaker: tokenMaker,
	}

	server.setupRoutes()

	return server
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	v1Routes := router.Group(server.basePath)

	v1Routes.GET("/heartbeat", server.getHeartbeat)
	v1Routes.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server.router = router
}

func (server *Server) Start() error {
	return server.router.Run(server.url)
}
