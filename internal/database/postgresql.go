package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Vladroon22/REST-API/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DataBase struct {
	logger *logrus.Logger
	config *config.Config
	sqlDB  *sqlx.DB
}

func NewDB(conf *config.Config, logg *logrus.Logger) *DataBase {
	return &DataBase{
		config: conf,
		logger: logg,
	}
}

func (d *DataBase) Connect() error {
	if err := d.openDB(*d.config); err != nil {
		d.logger.Errorln(err)
		return err
	}
	return nil
}

func (d *DataBase) openDB(conf config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	str := fmt.Sprintf("postgresql://%s", conf.DB)
	db, err := sqlx.ConnectContext(ctx, "postgres", str)
	d.logger.Infoln(str)
	if err != nil {
		d.logger.Errorln(err)
		return err
	}
	if err := RetryPing(db); err != nil {
		d.logger.Errorln(err)
		return err
	}
	d.sqlDB = db

	return nil
}

func RetryPing(db *sqlx.DB) error {
	var err error
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return err
}

func (db *DataBase) CloseDB() error {
	return db.sqlDB.Close()
}
