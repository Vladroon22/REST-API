package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Vladroon22/REST-API/config"
	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/handlers"
	server "github.com/Vladroon22/REST-API/internal/server"
	"github.com/Vladroon22/REST-API/internal/service"
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

	DB := db.NewDB(conf, logg)
	if err := DB.Connect(); err != nil {
		logg.Fatalln(err)
	}

	repo := db.NewRepo(DB)
	services := service.NewService(repo)
	router := handlers.NewRouter(services)

	srv := server.New(conf, logg)
	go func() {
		if err := srv.Run(router); err != nil || err != http.ErrServerClosed {
			logg.Fatalln(err)
		}
	}()

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM)

	<-killSig

	go func() {
		var wg sync.WaitGroup

		wg.Add(1)

		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := DB.CloseDB(); err != nil {
			logg.Errorln(err)
			return
		}
		if err := srv.Shutdown(ctx); err != nil {
			logg.Errorln(err)
			return
		}

		wg.Wait()

	}()

	logg.Infoln("Graceful shutdown...")

}
