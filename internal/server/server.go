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

func New(conf *config.Config, log *logrus.Logger) *Server {
	return &Server{
		server: &http.Server{},
		conf:   conf,
		logger: log,
	}
}

func (s *Server) Run(router *handlers.Router) error {
	s.logger.Infoln("Init router")

	router.Pref("/").SayHello()           // <-- logout
	router.Pref("/auth").AuthEndPoints()  // <-- sign-up/sign-in
	router.Pref("/users").UserEndPoints() // <-- only if sign-up/sign-in was success

	s.server = &http.Server{
		Addr:           s.conf.Addr_PORT,
		Handler:        router.R,
		MaxHeaderBytes: 1 << 20,
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
	}

	s.logger.Infoln("Server is listening -->", s.conf.Addr_PORT)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
