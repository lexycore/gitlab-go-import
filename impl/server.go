package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/xanzy/go-gitlab"
)

// Server is the main Server object structure
type Server struct {
	config *appConfig
	mux    *http.ServeMux
	web    *http.Server
	git    *gitlab.Client
}

// NewServer creates new Server object
func NewServer(config *appConfig) *Server {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(config.GitLabURL, "/") {
		config.GitLabURL += "/"
	}

	return &Server{
		config: config,
	}
}

// Serve prepares and starts Server
func (server *Server) Serve() error {
	if err := server.initGitLabClient(); err != nil {
		return err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	webErr := make(chan error, 1)

	server.setHandlers()

	go func() {
		listener, err := net.Listen("tcp", server.config.BindAddr)
		if err != nil {
			webErr <- err
			return
		}
		log.Printf("Serving at %s...", server.config.BindAddr)
		server.web = &http.Server{Handler: server.mux}
		if err = server.web.Serve(listener); err != nil {
			webErr <- err
			return
		}
		webErr <- nil
	}()

	var res error
	select {
	case err := <-webErr:
		if err != nil {
			log.Printf("got web server error: %v", err)
		}
		res = err
	case killSignal := <-interrupt:
		log.Printf("got signal: %v", killSignal)
		res = errors.New("server was interrupted by system signal")
	}

	if err := server.Shutdown(5 * time.Second); err != nil {
		log.Print(err)
	}
	return res
}

// Shutdown stops mux server
func (server *Server) Shutdown(timeout time.Duration) error {
	if server.web != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		log.Printf("stopping server...")
		defer cancel()
		return server.web.Shutdown(ctx)
	}
	return nil
}

func (server *Server) initGitLabClient() error {
	var err error
	server.git, err = gitlab.NewClient(server.config.GitLabToken, gitlab.WithBaseURL(fmt.Sprintf("%sapi/v4/", server.config.GitLabURL)))
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) setHandlers() {
	server.mux = http.NewServeMux()
	server.mux.HandleFunc("/go-get", server.goGet)
}
