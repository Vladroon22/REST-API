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

func (d *DataBase) ConfigDB() error {
	if err := d.openDB(*d.config); err != nil {
		d.logger.Errorln(err)
		return err
	}
	return nil
}

func (d *DataBase) openDB(conf config.Config) error {
	str := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBname, conf.SSLmode)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", str)
	d.logger.Infoln(str)
	if err != nil {
		d.logger.Errorln(err)
		return err
	}
	if err := db.Ping(); err != nil {
		d.logger.Errorln(err)
		return err
	}
	d.sqlDB = db
	d.logger.Infoln("Database configurated")

	return nil
}

func (db *DataBase) CloseDB() error {
	return db.sqlDB.Close()
}
