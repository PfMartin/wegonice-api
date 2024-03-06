package main

import (
	"log"

	"github.com/PfMartin/wegonice-api/config"
	"github.com/PfMartin/wegonice-api/db"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	_ = db.NewDatabase(conf.DBName, conf.DBUser, conf.DBPassword, conf.DBURI)
}
