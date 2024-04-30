package main

import (
	"fmt"

	"github.com/PfMartin/wegonice-api/api/v1"
	"github.com/PfMartin/wegonice-api/config"
	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/logging"
	"github.com/rs/zerolog/log"

	"github.com/PfMartin/wegonice-api/api/v1/docs"
)

func printBanner() {
	fmt.Print(`
██╗    ██╗███████╗ ██████╗  ██████╗ ███╗   ██╗██╗ ██████╗███████╗     █████╗ ██████╗ ██╗
██║    ██║██╔════╝██╔════╝ ██╔═══██╗████╗  ██║██║██╔════╝██╔════╝    ██╔══██╗██╔══██╗██║
██║ █╗ ██║█████╗  ██║  ███╗██║   ██║██╔██╗ ██║██║██║     █████╗█████╗███████║██████╔╝██║
██║███╗██║██╔══╝  ██║   ██║██║   ██║██║╚██╗██║██║██║     ██╔══╝╚════╝██╔══██║██╔═══╝ ██║
╚███╔███╔╝███████╗╚██████╔╝╚██████╔╝██║ ╚████║██║╚██████╗███████╗    ██║  ██║██║     ██║
 ╚══╝╚══╝ ╚══════╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝╚══════╝    ╚═╝  ╚═╝╚═╝     ╚═╝
																																												
`)
}

// @title											WeGoNice API
// @version										1.0
// @description 							API for the vegan recipes
// @termsOfService  					http://swagger.io/terms/

// @contact.name  						Martin Pfatrisch
// @contact.url   						https://github.com/PfMartin
// @contact.email 						martinpfatrisch@gmail.com

// @license.name  						All Rights Reserved

// @securityDefinitions.basic	BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @host											localhost:8000
// @BasePath									/api/v1
func main() {
	printBanner()

	logging.NewLogger()

	conf, err := config.NewConfig("./", ".env")
	if err != nil {
		log.Err(err).Msg("failed to read config")
		return
	}

	docs.SwaggerInfo.Title = "WeGoNice API"
	docs.SwaggerInfo.Description = "This is the WeGoNice API for managing vegan recipes and authors."
	docs.SwaggerInfo.Version = conf.APIVersion
	docs.SwaggerInfo.Host = conf.APIURL
	docs.SwaggerInfo.BasePath = conf.APIBasePath
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	wegoniceStore := db.NewMongoDBStore(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)

	server := api.NewServer(wegoniceStore, conf.APIURL, conf.APIBasePath, conf.TokenSymmetricKey, conf.AccessTokenDuration.Abs(), conf.RefreshTokenDuration.Abs())
	if err = server.Start(); err != nil {
		log.Err(err).Msg("failed to start server")
		return
	}
}
