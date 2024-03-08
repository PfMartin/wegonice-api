package main

import (
	"fmt"
	"log"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/PfMartin/wegonice-api/db"
	"github.com/PfMartin/wegonice-api/logging"
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
		log.Fatal(err)
	}

	_, cancel := db.NewDatabase(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
	defer cancel()
}
