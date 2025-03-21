package infra

import (
	"database/sql"
	"github.com/bubaew95/yandex-diploma/conf"
	"time"
)

type DataBase struct {
	*sql.DB
}

func NewDB(c *conf.Config) (*DataBase, error) {
	db, err := connectDB(c)
	if err != nil {
		return nil, err
	}

	return &DataBase{db}, nil
}

func connectDB(c *conf.Config) (*sql.DB, error) {
	db, err := sql.Open(c.Database.Driver, c.Database.DatabaseURI)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * time.Duration(c.Database.ConnMaxLifeTimeInMinute))
	db.SetMaxOpenConns(c.Database.MaxOpenConns)
	db.SetMaxIdleConns(c.Database.MaxIdleConns)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func (d DataBase) Close() error {
	return d.DB.Close()
}
