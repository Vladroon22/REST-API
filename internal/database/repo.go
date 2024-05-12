package database

import (
	"context"
	"time"
)

type repo struct {
	db      *DataBase
	timeOut time.Duration
}

func NewRepo(db *DataBase) *repo {
	return &repo{
		db:      db,
		timeOut: time.Duration(2) * time.Second,
	}
}

func (rp *repo) CreateNewUser(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	var id int
	if err := user.Valid(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	if err := user.HashingPass(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	query := "INSERT INTO clients (username, email, encrypt_password) VALUES ($1, $2, $3) RETURNING id"
	if err := rp.db.sqlDB.QueryRowContext(ctx, query, user.Name, user.Email, user.Encrypt_Password).Scan(&id); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully added")
	user.ID = id
	return id, nil
}

func (rp *repo) DeleteUser(c context.Context, id int) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "DELETE FROM clients WHERE id = $1 RETURNING id"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, id)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully deleted")
	user.ID = id
	return id, nil
}

func (rp *repo) UpdateUserFully(c context.Context, id int, name, email, pass string) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "UPDATE clients SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1 RETURNING id, username, email, encrypt_password"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, name, email, pass, id)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully updated")
	user.ID = id
	return id, nil
}

func (rp *repo) PartUpdateUserName(c context.Context, id int, name string) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "UPDATE clients SET username = $2 WHERE id = $1 RETURNING id, username"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, name, id)
	if err != nil {
		rp.db.logger.Infoln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' with his new name '%s'\n", id, name)
	return user.ID, nil
}

func (rp *repo) PartUpdateUserEmail(c context.Context, id int, email string) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "UPDATE users SET email = $3 WHERE id = $1 RETURNING id, email"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, email, id)
	if err != nil {
		rp.db.logger.Infoln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' and new email '%s'\n", id, email)
	return user.ID, nil
}

func (rp *repo) PartUpdateUserPass(c context.Context, id int, pass string) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "UPDATE users SET encrypt_password = $4 WHERE id = $1 RETURNING id, encrypt_password"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, pass, id)
	if err != nil {
		rp.db.logger.Infoln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' with his new password '%s'\n", id, pass)
	return user.ID, nil
}
