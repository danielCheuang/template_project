package redis

import "template_project/config"

var DB *Service

func Init() {
	var service = &Service{}
	cfg := config.GetConfig().Redis
	config := Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		Password:    cfg.Password,
		IdleTimeout: cfg.IdleTimeout,
		MaxActive:   cfg.MaxActive,
		MaxIdle:     cfg.MaxIdle,
	}
	service.Initialize(config)
	DB = service
}
