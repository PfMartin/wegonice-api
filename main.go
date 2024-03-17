package main

import (
	"fmt"

	"github.com/PfMartin/wegonice-api/api/v1"
	"github.com/PfMartin/wegonice-api/config"
	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/logging"
	"github.com/rs/zerolog/log"
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

func main() {
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
