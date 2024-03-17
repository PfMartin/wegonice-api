package api

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	dbClient *mongo.Client
	dbName   string
	url      string
	basePath string
	router   *gin.Engine
}

func NewServer(dbClient *mongo.Client, dbName string, url string, basePath string) *Server {
	server := &Server{
		dbClient: dbClient,
		dbName:   dbName,
		url:      url,
		basePath: basePath,
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