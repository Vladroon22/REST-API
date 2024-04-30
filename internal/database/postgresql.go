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

func NewDB(conf *config.Config) *DataBase {
	return &DataBase{
		config: conf,
		logger: logrus.New(),
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
	db, err := sql.Open("postgres", str)
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

func (db *DataBase) CloseDB() {
	db.sqlDB.Close()
}

func (db *DataBase) CreateNewUser(user *User) (int, error) {
	if err := user.Valid(); err != nil {
		db.logger.Errorln(err)
		return 0, err
	}
	if err := user.HashingPass(); err != nil {
		db.logger.Errorln(err)
		return 0, err
	}
	if err := db.sqlDB.QueryRow("INSERT INTO users (id, username, email, encrypt_password) VALUES ($1, $2, $3, $4) RETURNING id",
		user.ID, user.Name, user.Email, user.Encrypt_Password,
	).Scan(&user.ID); err != nil {
		db.logger.Errorln(err)
		return 0, err
	}

	db.logger.Infoln("User successfully added")
	return user.ID, nil
}

func (db *DataBase) DeleteUser(id int) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"DELETE FROM users WHERE id = $1 RETURNING id, username = $2, email = $3, encrypt_password = $4", id)
	if err != nil {
		db.logger.Errorln(err)
		return 0, err
	}

	db.logger.Infoln("User successfully deleted")
	return user.ID, nil
}

func (db *DataBase) UpdateUserFully(id int, name, email, pass string) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1 RETURNING id, username, email, encrypt_password", name, email, pass, id)
	if err != nil {
		db.logger.Errorln(err)
		return 0, err
	}

	db.logger.Infoln("User successfully updated")
	return user.ID, nil
}

func (db *DataBase) PartUpdateUserName(id int, name string) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET username = $2 WHERE id = $1 RETURNING id, username", name, id)
	if err != nil {
		db.logger.Infoln(err)
		return 0, err
	}

	db.logger.Infof("Update the User '%d' with his new name '%s'\n", id, name)
	return user.ID, nil
}

func (db *DataBase) PartUpdateUserEmail(id int, email string) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET email = $3 WHERE id = $1 RETURNING id, email", email, id)
	if err != nil {
		db.logger.Infoln(err)
		return 0, err
	}

	db.logger.Infof("Update the User '%d' and new email '%s'\n", id, email)
	return user.ID, nil
}

func (db *DataBase) PartUpdateUserPass(id int, pass string) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET encrypt_password = $4 WHERE id = $1 RETURNING id, encrypt_password", pass, id)
	if err != nil {
		db.logger.Infoln(err)
		return nil, err
	}

	db.logger.Infof("Update the User '%d' with his new password '%s'\n", id, pass)
	return user, nil
}
