package httpserver

import (
	"net"
	"time"
)

type OptionServer func(*Server)

func Port(port string) OptionServer {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) OptionServer {
	return func(s *Server) {
		s.server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) OptionServer {
	return func(s *Server) {
		s.server.WriteTimeout = timeout
	}
}

func ShutdownTimeout(timeout time.Duration) OptionServer {
	return func(s *Server) {
		s.shutdownTimeout = timeout
	}
}
