package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type heartBeatResponse struct {
	Status string `json:"status"`
}

func (server *Server) getHeartbeat(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, heartBeatResponse{
		Status: "ok",
	})
}
