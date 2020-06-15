package db

import (
	"fmt"
	"github.com/CodFrm/iotqq-plugins/config"
	goRedis "github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
)

var Redis *goRedis.Client
var Db *gorm.DB

func Init() error {
	Redis = goRedis.NewClient(&goRedis.Options{
		Addr:     config.AppConfig.Redis.Addr,
		Password: config.AppConfig.Redis.Password,
		DB:       config.AppConfig.Redis.DB,
	})
	if _, err := Redis.Ping().Result(); err != nil {
		return fmt.Errorf("redis open error: %v", err)
	}
	var err error
	Db, err = gorm.Open("mysql", config.AppConfig.MySQL.Dsn)
	if err != nil {
		return fmt.Errorf("sql open error: %v", err)
	}
	Db.SingularTable(true)
	Db.DB().SetMaxOpenConns(90)
	Db.DB().SetMaxIdleConns(50)

	driver, err := mysql.WithInstance(Db.DB(), &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"mysql", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
