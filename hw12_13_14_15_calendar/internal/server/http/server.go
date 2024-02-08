package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/spmadness/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	logger  Logger
	app     Application
	server  *http.Server
	address string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface {
	GetStorage() app.Storager
}

func NewServer(logger Logger, app Application, address string) *Server {
	return &Server{
		logger:  logger,
		app:     app,
		address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		mux := NewMux(s.app.GetStorage())

		mw := Middleware{
			logger: s.logger,
		}

		chain := MiddlewareChain(
			mw.loggingMiddleware,
			mw.requestValidatorMiddleware,
		)(mux)

		s.server = &http.Server{
			Addr:              s.address,
			Handler:           chain,
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		}

		s.logger.Info(fmt.Sprintf("starting http server on http://%s", s.address))
		err := s.server.ListenAndServe()
		if err != nil {
			s.logger.Error(fmt.Sprintf("server error: %s", err))
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("shutting down http server...")
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	s.logger.Info("http server shutdown success")
	return nil
}
