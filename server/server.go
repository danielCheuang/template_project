package server

import (
	"errors"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	Server *http.Server
	G      *errgroup.Group
	*API
}

func GetServer(api *API) (server *Server) {
	serverConfig := api.config.Server
	runMode := api.config.Server.RunMode

	currentEngineConfig := &EngineConfig{
		middleware:       nil,
		LimitConnections: serverConfig.LimitConnection,
		RunMode:          runMode,
	}

	return &Server{
		G:   api.ErrorGroup,
		API: api,
		Server: &http.Server{
			Handler:        currentEngineConfig.Init(api),
			ReadTimeout:    serverConfig.ReadTimeout,
			WriteTimeout:   serverConfig.WriteTimeout,
			IdleTimeout:    serverConfig.IdleTimeout,
			MaxHeaderBytes: serverConfig.MaxHeaderBytes,
		},
	}
}

func (server *Server) Run() error {

	server.runServer()

	if server.config.Server.EnableHTTPS {
		if server.config.TLS.CertFile == "" || server.config.TLS.KeyFile == "" {
			return errors.New("use https should config the cert and key files")
		}
		server.runServerTLS()
	}
	if err := server.G.Wait(); err != nil {
		return err
	}
	return nil
}

func (server *Server) runServer() {
	server.G.Go(func() error {
		return http.ListenAndServe(server.config.Server.ListenAddr, server.Server.Handler)
	})
}

func (server *Server) runServerTLS() {
	server.G.Go(func() error {
		return http.ListenAndServeTLS(server.config.Server.HTTPSAddr, server.config.TLS.CertFile, server.config.TLS.KeyFile, server.Server.Handler)
	})
}
