package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wallet-api/pkg/logger"

	"github.com/rs/zerolog"
)

type Server interface {
	Serve(parentContext context.Context)
	Router
}

type server struct {
	logger zerolog.Logger
	mux    *http.ServeMux
	*http.Server
	config ServerConfig
}

func NewServer(log zerolog.Logger, cfg ServerConfig) Server {
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		Handler:           mux,
	}
	log = logger.WithModule(log, "server")
	return &server{logger: log, mux: mux, Server: srv, config: cfg}
}

func (s *server) Serve(parentContext context.Context) {
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatal().Err(err).Msg("server fault")
		}
	}()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.logger.Info().Msg("shutting down server...")
	ctx, cancel := context.WithTimeout(parentContext, s.config.ShutdownTimeout)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		s.logger.Fatal().Err(err).Msg("server shutdown")
	}
	<-ctx.Done()
	s.logger.Info().Msg("server exiting")
}
