package server

import (
	"context"
	"net/http"
	"time"

	"github.com/Vladroon22/REST-API/config"
	"github.com/Vladroon22/REST-API/internal/handlers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	conf   *config.Config
	logger *logrus.Logger
	server *http.Server
}

func New(conf *config.Config, log *logrus.Logger, serv *http.Server) *Server {
	return &Server{
		server: serv,
		conf:   conf,
		logger: log,
	}
}

func (s *Server) Run() {
	s.logger.Infof("Listening: '%s'\n", s.conf.Addr_PORT)

	router := handlers.NewRouter()
	s.logger.Infoln("Created New router")

	router.Pref("/").SayHello()                // <-- logout
	router.Pref("/auth").EndPoints()           // <-- sign-up/sign-in
	router.Pref("/auth/users").UserEndPoints() // only if reg was success

	s.server = &http.Server{
		Addr:         s.conf.Addr_PORT,
		Handler:      router.R,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Infoln("Server is listening -->")
	go s.logger.Fatalln(s.server.ListenAndServe())
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
