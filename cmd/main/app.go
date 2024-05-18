package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/Vladroon22/REST-API/config"
	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/handlers"
	"github.com/Vladroon22/REST-API/internal/server"

	"github.com/sirupsen/logrus"
)

var (
	pathToToml string
)

func main() {
	flag.Parse()

	flag.StringVar(&pathToToml, "path-to-toml", "./config/conf.toml", "path-to-toml")

	logg := logrus.New()          // logger
	conf := config.CreateConfig() // config

	_, err := toml.DecodeFile(pathToToml, conf)
	if err != nil {
		logg.Errorln(err)
		return
	}

	DB := db.NewDB(conf)
	if err := DB.ConfigDB(); err != nil {
		logg.Fatalln(err)
	}

	repo := db.NewRepo(DB)
	router := handlers.NewRouter(repo)

	server.New(conf, logg).Run(router)
}
