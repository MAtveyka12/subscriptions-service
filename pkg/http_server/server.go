package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

type HTTPServer struct {
	logger *slog.Logger
	server *http.Server
	config *Config
}

type Config struct {
	Host              string
	Port              string
	StartMsg          string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

func NewServer(logger *slog.Logger, config *Config) *HTTPServer {
	server := &http.Server{
		Handler:           config.Handler,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		Addr:              config.Host + ":" + config.Port,
	}

	s := &HTTPServer{
		logger: logger,
		server: server,
		config: config,
	}

	return s
}

func (s *HTTPServer) Start(ctx context.Context) error {
	s.logger.Info(s.config.StartMsg)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			s.logger.Info("Context cancelled, shutting down gracefully...")
		case <-sigCh:
			s.logger.Info("Shutdown signal received, shutting down gracefully...")
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error("Failed to shutdown the server", slog.String("error", err.Error()))
			return err
		}

		s.logger.Info("Server is shutdown!")

		return nil
	})

	g.Go(func() error {
		err := s.server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			return err
		}

		return nil
	})

	return g.Wait()
}
