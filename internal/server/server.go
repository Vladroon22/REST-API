package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vladroon22/REST-API/config"
	"github.com/Vladroon22/REST-API/internal/database"
	"github.com/sirupsen/logrus"
)

type Server struct {
	conf   *config.Config
	logger *logrus.Logger
	server *http.Server
	db     *database.DataBase
}

func New(conf *config.Config, log *logrus.Logger) *Server {
	return &Server{
		server: &http.Server{},
		conf:   conf,
		logger: log,
		db:     database.NewDB(conf),
	}
}

func (s *Server) Run() {
	if err := s.db.ConfigDB(); err != nil {
		s.logger.Fatalln(err)
	}

	s.logger.Infof("Listening: '%s'\n", s.conf.Addr_PORT)

	router := database.NewRouter(s.db)
	s.logger.Infoln("Created New router")

	router.Pref("/").SayHello()           // <-- logout
	router.Pref("/auth").AuthEndPoints()  // <-- sign-up/sign-in
	router.Pref("/users").UserEndPoints() // <-- only if sign-up/sign-in was success

	s.server = &http.Server{
		Addr:         s.conf.Addr_PORT,
		Handler:      &router.R,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		s.logger.Infoln("Server is listening -->")
		s.logger.Fatalln(s.server.ListenAndServe())
	}()

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM)

	<-killSig

	go func() {
		if err := s.Shutdown(context.Background()); err != nil {
			s.logger.Fatalln(err)
		}
	}()

}

func (s *Server) Shutdown(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, 3*time.Second)
	defer cancel()
	if err := s.db.CloseDB(); err != nil {
		s.logger.Errorln(err)
		return err
	}
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Errorln(err)
		return err
	}
	return nil
}
