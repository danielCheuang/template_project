package server

import (
	"template_project/config"

	"golang.org/x/sync/errgroup"
)

var eGroup errgroup.Group

type API struct {
	config     *config.Configuration
	ErrorGroup *errgroup.Group
}

func New(conf *config.Configuration) (*API, error) {
	cfg := *conf
	return &API{
		config:     &cfg,
		ErrorGroup: &eGroup,
	}, nil
}

func (api *API) Start() error {
	server := GetServer(api)
	return server.Run()
}
