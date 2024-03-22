package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type heartbeatResponse struct {
	Status string `json:"status" example:"ok"`
} // @name heartbeatResponse

// getHeartbeat
// @Summary 			Check heartbeat
// @Schemes 			http,https
// @Description 	Check if the API is reachable with this route
// @Tags 					heartbeat
// @Accept 				application/json
// @Produce 			application/json
// @Success 			200 				{object} 	heartbeatResponse 		"Success"
// @Router 				/heartbeat 	[get]
func (server *Server) getHeartbeat(ctx *gin.Context) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, heartbeatResponse{
		Status: "ok",
	})
}
