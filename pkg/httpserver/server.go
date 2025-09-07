package httpserver

import (
	"context"
	"net/http"
	"time"
)

const _defaultShutdownTimeout time.Duration = 5 * time.Second

type Server struct {
	instance        *http.Server
	notify          chan error
	isHttps         bool
	certFile        string
	keyFile         string
	shutdownTimeout time.Duration
}

func NewServer(router http.Handler, address string, isHttps bool, certFile, keyFile string) *Server {
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &Server{
		instance:        server,
		notify:          make(chan error, 1),
		isHttps:         isHttps,
		certFile:        certFile,
		keyFile:         keyFile,
		shutdownTimeout: _defaultShutdownTimeout,
	}
}

func (s *Server) Start() {
	if s.isHttps {
		go func() {
			s.notify <- s.instance.ListenAndServeTLS(s.certFile, s.keyFile)

			close(s.notify)
		}()
	} else {
		go func() {
			s.notify <- s.instance.ListenAndServe()

			close(s.notify)
		}()
	}
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := s.instance.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
