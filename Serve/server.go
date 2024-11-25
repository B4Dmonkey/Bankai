package serve

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ServerConfigFunc func(*ServerConfig)

type ServerConfig struct {
	ctx             context.Context
	cancel          context.CancelFunc
	mux             *http.ServeMux
	server          *http.Server
	shutdownTimeout time.Duration
}

type Server struct {
	ServerConfig
}

func defaultServerConfig(addr string) (*ServerConfig, error) {
	ctx, cancel := context.WithCancel(context.Background())
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return &ServerConfig{ctx: ctx, cancel: cancel, mux: mux, server: server, shutdownTimeout: 5 * time.Second}, nil
}

func New(overrides ...ServerConfigFunc) *Server {
	var errs []error

	config, err := defaultServerConfig("")
	if err != nil {
		errs = append(errs, err)
	}

	for _, override := range overrides {
		override(config)
	}

	return &Server{
		ServerConfig: *config,
	}
}

func (s *Server) Start(wg *sync.WaitGroup) {
	log.Println("Starting server...")
	defer wg.Done()
	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Listening and serving at", s.server.Addr)
		log.Fatal(s.server.ListenAndServe(), "Server failed to start")
	}()

	<-gracefulShutdown
	println() // * Formatting for log readability
	log.Println("Received terminate, gracefully shutting down...")

	_, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	log.Println("Server stopped")
}
