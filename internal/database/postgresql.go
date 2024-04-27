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

/*
	type Storage interface {
		CreateNewUser(user *User) (*User, error)
		DeleteUser(id int) (*User, error)
		UpdateUserFully(id int, name, email, pass string) (*User, error)
		PartUpdateUserName(id int, name string) (*User, error)
		PartUpdateUserEmail(id int, email string) (*User, error)
		PartUpdateUserPass(id int, pass string) (*User, error)
	}
*/
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

func (db *DataBase) CreateNewUser(user *User) (*User, error) {
	if err := user.Valid(); err != nil {
		db.logger.Infoln(err)
		return nil, err
	}
	if err := user.HashingPass(); err != nil {
		db.logger.Errorln(err)
		return nil, err
	}
	if err := db.sqlDB.QueryRow("INSERT INTO users (username, email, encrypt_password) VALUES ($2, $3, $4) RETURNING id",
		user.ID, user.Name, user.Email, user.Encrypt_Password,
	).Scan(&user.ID); err != nil {
		db.logger.Errorln(err)
		return nil, err
	}

	db.logger.Infoln("User successfully added")
	return user, nil
}

func (db *DataBase) DeleteUser(id int) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"DELETE FROM users WHERE id = $1 RETURNING id, username, email, encrypt_password", id)
	if err != nil {
		db.logger.Errorln(err)
		return nil, err
	}

	db.logger.Infoln("User successfully deleted")
	return user, nil
}

func (db *DataBase) UpdateUserFully(id int, name, email, pass string) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1 RETURNING id, username, email, encrypt_password", name, email, pass, id)
	if err != nil {
		db.logger.Errorln(err)
		return nil, err
	}

	db.logger.Infoln("User successfully updated")
	return user, nil
}

func (db *DataBase) PartUpdateUserName(id int, name string) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET username = $2 WHERE id = $1 RETURNING id, username", id, name)
	if err != nil {
		db.logger.Infoln(err)
		return nil, err
	}

	db.logger.Infof("Update the User '%d' with his new name '%s'\n", id, name)
	return user, nil
}

func (db *DataBase) PartUpdateUserEmail(id int, email string) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET email = $3 WHERE id = $1 RETURNING id, email", id, email)
	if err != nil {
		db.logger.Infoln(err)
		return nil, err
	}

	db.logger.Infof("Update the User '%d' and new email '%s'\n", id, email)
	return user, nil
}

func (db *DataBase) PartUpdateUserPass(id int, pass string) (*User, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET username = $4 WHERE id = $1 RETURNING id, encrypt_password", id, pass,
	)
	if err != nil {
		db.logger.Infoln(err)
		return nil, err
	}

	db.logger.Infof("Update the User '%d' with his new password '%s'\n", id, pass)
	return user, nil
}
