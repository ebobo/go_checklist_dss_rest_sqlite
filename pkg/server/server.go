package server

import (
	"context"
	"fmt"
	"sync"
)

// Server takes care of instantiating and running service and other dependencies.
type Server struct {
	httpListenAddr string
	httpStarted    *sync.WaitGroup
	httpStopped    *sync.WaitGroup
	ctx            context.Context
	cancel         context.CancelFunc
}

// Config is the server configuration
type Config struct {
	GRPCListenAddr string
	HTTPListenAddr string
	MSGPRCAddr     string
}

func New(c Config) *Server {
	return &Server{
		httpListenAddr: c.HTTPListenAddr,
		httpStarted:    &sync.WaitGroup{},
		httpStopped:    &sync.WaitGroup{},
	}
}

func (s *Server) Start() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// Start the HTTP interface
	s.httpStarted.Add(1)
	s.httpStopped.Add(1)
	err := s.startHTTP()
	if err != nil {
		return err
	}
	s.httpStarted.Wait()

	return nil
}

func (s *Server) Shutdown() {
	fmt.Println("shut down")
	if s.cancel != nil {
		s.cancel()
	}
	s.httpStopped.Wait()
}
