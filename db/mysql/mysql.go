package mysql

import (
	"template_project/config"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *Service

func Init() {

	cfg := config.GetConfig().MySQL
	config := Config{
		Host:        cfg.Host,
		Port:        cfg.Port,
		User:        cfg.User,
		Password:    cfg.Password,
		DbName:      cfg.DbName,
		TablePrefix: cfg.TablePrefix,

		MaxOpenConnections: cfg.MaxOpenConnections,
		MaxIdleConnections: cfg.MaxIdleConnections,
		ConnMaxLifetime:    cfg.ConnMaxLifetime,
		Debug:              cfg.Debug,
		Local:              cfg.Local,
	}
	db, err := NewService(config)
	if err != nil {
		panic(err)
	}
	DB = db
}
