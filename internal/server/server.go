package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Vladroon22/REST-API/config"
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

func (s *Server) Run(router handlers.) {
	s.logger.Infof("Listening: '%s'\n", s.conf.Addr_PORT)

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
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalln(err)
		}
	}()

	killSig := make(chan os.Signal, 1)
	signal.Notify(killSig, syscall.SIGINT, syscall.SIGTERM)

	<-killSig

	go func() {
		if err := s.Shutdown(context.Background()); err != nil {
			s.logger.Fatalln(err)
		}
	}()
	s.logger.Infoln("Graceful shutdown...")
}

func (s *Server) Shutdown(c context.Context) error {
	var wg sync.WaitGroup

	wg.Add(1)

	defer wg.Done()

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	if err := s.db.CloseDB(); err != nil {
		return err
	}
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	wg.Wait()

	return nil
}
