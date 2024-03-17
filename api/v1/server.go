package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	dbClient *mongo.Client
	dbName   string
	url      string
	router   *gin.Engine
}

func NewServer(dbClient *mongo.Client, dbName string, url string) *Server {
	server := &Server{
		dbClient: dbClient,
		dbName:   dbName,
		url:      url,
	}

	server.setupRoutes()

	return server
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	v1Routes := router.Group("/api/v1")

	v1Routes.GET("/heartbeat", server.getHeartbeat)

	server.router = router
}

func (server *Server) Start() error {
	return server.router.Run(server.url)
}
