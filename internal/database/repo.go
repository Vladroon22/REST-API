package database

import (
	"context"
)

type repo struct {
	db *DataBase
}

func NewRepo(db *DataBase) *repo {
	return &repo{
		db: db,
	}
}

func (rp *repo) CreateNewUser(ctx context.Context, user *User) (int, error) {
	var id int
	if err := user.Valid(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	if err := user.HashingPass(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	query := "INSERT INTO clients (id, username, email, encrypt_password) VALUES ($1, $2, $3, $4) RETURNING id"
	if err := rp.db.sqlDB.QueryRowContext(ctx, query, user.ID, user.Name, user.Email, user.Encrypt_Password).Scan(&id); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully added")
	user.ID = id
	return id, nil
}

func (rp *repo) DeleteUser(id int) (int, error) {
	user := &User{}
	_, err := rp.db.sqlDB.Exec(
		"DELETE FROM clients WHERE id = $1 RETURNING id, username = $2, email = $3, encrypt_password = $4", id)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully deleted")
	return user.ID, nil
}

func (rp *repo) UpdateUserFully(id int, name, email, pass string) (int, error) {
	user := &User{}
	_, err := rp.db.sqlDB.Exec(
		"UPDATE clients SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1 RETURNING id, username, email, encrypt_password", name, email, pass, id)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully updated")
	return user.ID, nil
}

func (db *DataBase) PartUpdateUserName(id int, name string) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE clients SET username = $2 WHERE id = $1 RETURNING id, username", name, id)
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

func (db *DataBase) PartUpdateUserPass(id int, pass string) (int, error) {
	user := &User{}
	_, err := db.sqlDB.Exec(
		"UPDATE users SET encrypt_password = $4 WHERE id = $1 RETURNING id, encrypt_password", pass, id)
	if err != nil {
		db.logger.Infoln(err)
		return 0, err
	}

	db.logger.Infof("Update the User '%d' with his new password '%s'\n", id, pass)
	return user.ID, nil
}
