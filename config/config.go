package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/configor"
)

var (
	Cfg  Configuration
	Lock sync.RWMutex
)

type (
	ServerConfig struct {
		Name             string        `json:"name"`
		RunMode          string        `json:"run_mode"`
		ListenAddr       string        `json:"listen_addr"`
		LimitConnection  int           `json:"limit_connection"`
		RootRouterPrefix string        `json:"root_router_prefix"`
		EnableHTTPS      bool          `json:"enable_https"`
		HTTPSAddr        string        `json:"https_addr"`
		ReadTimeout      time.Duration `json:"read_timeout"`
		WriteTimeout     time.Duration `json:"write_timeout"`
		IdleTimeout      time.Duration `json:"idle_timeout"`
		MaxHeaderBytes   int           `json:"max_header_bytes"`
	}

	TLSConfig struct {
		CertFile string `json:"cert_file"`
		KeyFile  string `json:"key_file"`
	}

	MySQLConfig struct {
		Enable             bool   `json:"enable" defualt:"false"`
		Host               string `json:"host"`
		Port               string `json:"port"`
		User               string `json:"user"`
		Password           string `json:"password"`
		DbName             string `json:"db_name"`
		TablePrefix        string `json:"table_prefix"`
		MaxOpenConnections int    `json:"max_open_connections"`
		MaxIdleConnections int    `json:"max_idle_connections"`
		ConnMaxLifetime    int    `json:"conn_max_lifetime"`
		Local              string `json:"local"`
		Debug              bool   `json:"debug"`
	}

	RedisConfig struct {
		Enable      bool   `json:"enable" defualt:"false"`
		Host        string `json:"host"`
		Port        string `json:"port"`
		Password    string `json:"password"`
		IdleTimeout int    `json:"idle_timeout"`
		MaxIdle     int    `json:"max_idle"`
		MaxActive   int    `json:"max_active"`
	}

	LoggerConfig struct {
		Level          string        `json:"level"`
		Formatter      string        `json:"formatter"`
		DisableConsole bool          `json:"disable_console"`
		Write          bool          `json:"write"`
		Path           string        `json:"path"`
		FileName       string        `json:"file_name"`
		MaxAge         time.Duration `json:"max_age"`
		RotationTime   time.Duration `json:"rotation_time"`
		Debug          bool          `json:"debug"`
	}

	ChainConfig struct {
		Account string `json:"account"`
		Secret  string `json:"secret"`
	}

	Configuration struct {
		Server ServerConfig           `json:"server"`
		TLS    TLSConfig              `json:"tls"`
		MySQL  MySQLConfig            `json:"mysql"`
		Redis  RedisConfig            `json:"redis"`
		Logger LoggerConfig           `json:"logger"`
		Chains map[string]ChainConfig `json:"chains"`
	}
)

func Init(file *string) Configuration {
	Lock.Lock()
	defer Lock.Unlock()

	err := configor.Load(&Cfg, *file)
	if err != nil {
		fmt.Println(err)
	}
	return Cfg
}

func GetConfig() Configuration {
	Lock.Lock()
	defer Lock.Unlock()
	return Cfg
}
