package database

import (
	"database/sql"
	"fmt"

	"github.com/Vladroon22/REST-API/config"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type DataBase struct {
	logger *logrus.Logger
	config *config.Config
	sqlDB  *sql.DB
}

func NewDB(conf *config.Config, log *logrus.Logger) *DataBase {
	return &DataBase{
		config: conf,
		logger: log,
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
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBname))
	if err != nil {
		d.logger.Errorln(err)
		return err
	}

	if err := db.Ping(); err != nil {
		d.logger.Errorln("Wrong config: ", err)
	}
	d.sqlDB = db
	d.logger.Infoln("Database configurated")

	return nil
}

func (db *DataBase) CloseDB() error {

	return nil
}