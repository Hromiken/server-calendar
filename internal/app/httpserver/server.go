package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
)

// Server — обёртка над http.Server с поддержкой graceful shutdown
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New — создаёт и сразу запускает HTTP-сервер
func New(handler http.Handler, opts ...OptionServer) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		Addr:         defaultAddr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

// start — запускает сервер в отдельной горутине
func (s *Server) start() {
	go func() {
		logrus.Infof("🚀 Starting HTTP server on %s", s.server.Addr)
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.notify <- err
			logrus.Errorf("❌ Server stopped with error: %v", err)
		} else {
			logrus.Info("🟢 HTTP server stopped gracefully")
		}
		close(s.notify)
	}()

}

// Notify — канал для ошибок сервера
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown — корректное завершение работы
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	logrus.Infof("🛑 Shutting down HTTP server (timeout %v)...", s.shutdownTimeout)

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	logrus.Info("✅ Server shutdown complete")
	return nil
}
