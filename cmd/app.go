package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Vladroon22/REST-API/config"
	d "github.com/Vladroon22/REST-API/internal/database"
	"github.com/Vladroon22/REST-API/internal/handlers"
	server "github.com/Vladroon22/REST-API/internal/server"
	"github.com/Vladroon22/REST-API/internal/service"
	"github.com/sirupsen/logrus"
)

// @title REST-API
// @version 1.0
// @description API

// @host 127.0.0.1:8000
// @BasePath /

// @securityDefinitions.apikey signKey
// @in header
// @name jwt

func main() {
	logg := logrus.New()
	conf := config.CreateConfig()

	db := d.NewDB(conf, logg)
	if err := db.Connect(); err != nil {
		logg.Fatalln(err)
	}

	repo := d.NewRepo(db)
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

		if err := db.CloseDB(); err != nil {
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
