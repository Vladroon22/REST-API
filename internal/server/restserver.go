package server

import (
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

	router.Pref("/").SayHello()
	router.Pref("/auth").EndPoints()
	router.Pref("/users").UserEndPoints()

	s.server = &http.Server{
		Addr:         s.conf.Addr_PORT,
		Handler:      router.R,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.logger.Infoln("Server is listening -->")
	s.logger.Fatalln(s.server.ListenAndServe())
}

/*
func shutdown(s *Server, ctx context.Context) {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGINT)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	s.HttpServer.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	s.Logger.Infoln("Graceful shutting down")
	os.Exit(0)
}
*/
