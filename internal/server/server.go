package server

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	s *http.Server
}

func New(port string, h http.Handler) Server {
	return Server{
		s: &http.Server{
			Addr:         net.JoinHostPort("0.0.0.0", port),
			Handler:      h,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
	}
}

// ListenAndServe listens for connections on interface 0.0.0.0 at the port
// provided to New. If the provided context is canceled, the server will attempt
// to gracefully shutdown. The returned error will only be non-nil if the server
// exits abnormally or fails to shutdown in time.
func (s Server) ListenAndServe(ctx context.Context) error {
	errC := make(chan error, 1)
	go func() { errC <- s.s.ListenAndServe() }()

	select {
	case <-ctx.Done():
		return s.shutdown()
	case err := <-errC:
		if err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

func (s Server) shutdown() error {
	// We have to create a new context here rather than passing one in from
	// ListenAndServe since that one was already Done'd by the time this method is
	// called.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.s.Shutdown(ctx)
}
