package logger

import (
	"template_project/config"

	"github.com/sirupsen/logrus"
)

var Log Logger

type Logger struct {
	log *logrus.Logger
}

func Init() {
	cfg := config.GetConfig().Logger
	Log.log = NewService(&Config{
		Level:          cfg.Level,
		Formatter:      cfg.Formatter,
		DisableConsole: cfg.DisableConsole,
		Write:          cfg.Write,
		Path:           cfg.Path,
		FileName:       cfg.FileName,
		MaxAge:         cfg.MaxAge,
		RotationTime:   cfg.RotationTime,
		Debug:          cfg.Debug,
	})
}
