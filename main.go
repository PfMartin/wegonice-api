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

// // @termsOfService  http://swagger.io/terms/

// @contact.name   Martin Pfatrisch
// @contact.url    https://github.com/PfMartin
// @contact.email  martinpfatrisch@gmail.com

// @license.name  All Rights Reserved
// // @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	docs.SwaggerInfo.Title = "WeGoNice API"
	docs.SwaggerInfo.Description = "This is the WeGoNice API for vegan recipes."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	printBanner()

	logging.NewLogger()

	conf, err := config.NewConfig("./", ".env")
	if err != nil {
		log.Err(err).Msg("failed to read config")
		return
	}

	dbClient, cancel := db.NewDatabaseClient(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	_ = db.NewUserCollection(dbClient, conf.DBName)
	_ = db.NewAuthorCollection(dbClient, conf.DBName)

	server := api.NewServer(dbClient, conf.DBName, conf.APIURL)
	if err = server.Start(); err != nil {
		log.Err(err).Msg("failed to start server")
		return
	}

	defer cancel()
}
