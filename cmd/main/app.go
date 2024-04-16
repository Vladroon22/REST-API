package main

import (
	"flag"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/Vladroon22/REST-API/config"
	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/server"
	"github.com/sirupsen/logrus"
)

var (
	srv        *http.Server
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

	d := db.NewDB(conf, logg)

	if err := d.ConfigDB(); err != nil {
		logg.Fatalln(err)
	}

	server.New(conf, logg, srv).Run()
}
