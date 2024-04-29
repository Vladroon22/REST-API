package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/Vladroon22/REST-API/config"
	db "github.com/Vladroon22/REST-API/internal/database"
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

	d := db.NewDB(conf)

	server.New(conf, logg, d).Run()
}
