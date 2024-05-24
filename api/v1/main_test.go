package api

import (
	"os"
	"testing"
	"time"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/util"
	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}

func newTestServer(t *testing.T, store db.DBStore) *Server {
	t.Helper()

	config := config.Config{
		TokenSymmetricKey:    util.RandomString(32),
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Minute,
		APIURL:               "localhost:8001",
		APIBasePath:          "/api/v1",
		CorsAllowedOrigins:   []string{"http://*", "https://*"},
		ImagesDepotPath:      "../../images/test_images_depot",
	}

	server := NewServer(
		store,
		config.APIURL,
		config.APIBasePath,
		config.TokenSymmetricKey,
		config.AccessTokenDuration,
		config.RefreshTokenDuration,
		config.CorsAllowedOrigins,
		config.ImagesDepotPath,
	)

	return server
}
