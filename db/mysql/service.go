package mysql

import (
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type Config struct {
	Host        string
	Port        string
	User        string
	Password    string
	DbName      string
	TablePrefix string

	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    int // unit second
	Debug              bool
	Local              string
}

type Service struct {
	*gorm.DB
	config Config
}

func NewService(config Config) (*Service, error) {
	impl := &Service{}
	impl.config = config
	var url string

	password := config.Password
	if password != "" {
		password = fmt.Sprintf(":%s", password)
	}

	url = fmt.Sprintf("%s%s@%s(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=%s", config.User, password, "tcp", config.Host, config.Port,
		config.DbName, config.Local)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.TablePrefix + defaultTableName
	}
	db, err := gorm.Open("mysql", url)
	if err != nil {
		log.Fatalf("mysql connection error:%s", err.Error())
		return nil, err
	}

	maxOpenConns := config.MaxOpenConnections
	if maxOpenConns < 5 {
		maxOpenConns = 5
	}

	db.DB().SetMaxOpenConns(maxOpenConns)

	maxIdleConns := config.MaxIdleConnections
	if maxIdleConns < 1 {
		maxIdleConns = 1
	}
	db.DB().SetMaxIdleConns(maxIdleConns)

	connMaxLifeTime := config.ConnMaxLifetime
	if connMaxLifeTime < 30 {
		connMaxLifeTime = 30
	}
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifeTime) * time.Second)

	db.LogMode(config.Debug)
	err = db.DB().Ping()
	if err != nil {
		panic(fmt.Sprintf("init mysql db err: %v", err))
	}

	impl.DB = db

	return impl, nil
}

// RegistTable create table for given object
func (s *Service) RegistTable(t interface{}) {
	if ok := s.DB.HasTable(t); !ok {
		if err := s.DB.CreateTable(t).Error; err != nil {
			log.Fatalf("create mysql table error:%s", err.Error())
		}
	}
	var tab []interface{}
	s.DB.AutoMigrate(append(tab, t))
}

// RegistTables create tables for given object
func (s *Service) RegistTables(tables []interface{}) {
	for _, t := range tables {
		if ok := s.DB.HasTable(t); !ok {
			if err := s.DB.CreateTable(t).Error; err != nil {
				log.Fatalf("create mysql table error:%s", err.Error())
			}
		}
	}

	// auto migrate to keep schema update to date
	// AutoMigrate will ONLY create tables, missing columns and missing indexes,
	// and WON'T change existing column's type or delete unused columns to protect your data
	s.DB.AutoMigrate(tables...)
}

// Close db
func (s *Service) Close() error {
	return s.DB.Close()
}

// Add single item
func (s *Service) Add(item interface{}) error {
	return s.DB.Create(item).Error
}

// Del single item
func (s *Service) Del(item interface{}) error {
	return s.DB.Delete(item).Error
}

// Save single item
func (s *Service) Save(item interface{}) error {
	return s.DB.Save(item).Error
}
