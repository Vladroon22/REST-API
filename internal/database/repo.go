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

	if ok, err := rp.IdExist(ctx, id); ok {
		return 0, err
	}

	query := "DELETE FROM clients WHERE id = $1 RETURNING id"
	rows, err := rp.db.sqlDB.ExecContext(ctx, query, id)
	res, _ := rows.RowsAffected()
	if res == 0 {
		rp.db.logger.Infoln("Database is empty")
		return 0, nil
	}
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully deleted")
	return id, nil
}

func (rp *repo) UpdateUserFully(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if ok, err := rp.IdExist(ctx, user.ID); ok {
		return 0, err
	}

	if err := user.Valid(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	if err := user.HashingPass(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	query := "UPDATE clients SET username = $2, email = $3, encrypt_password = $4 WHERE id = $1"

	stmt, err := rp.db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.ExecContext(ctx, user.ID, user.Name, user.Email, user.Encrypt_Password)
	res, _ := rows.RowsAffected()
	if res == 0 {
		rp.db.logger.Infoln("Database is empty")
		return 0, nil
	}
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully updated")
	return user.ID, nil
}

func (rp *repo) PartUpdateUserName(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if ok, err := rp.IdExist(ctx, user.ID); ok {
		return 0, err
	}

	query := "UPDATE clients SET username = $2 WHERE id = $1 RETURNING id"

	stmt, err := rp.db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.ExecContext(ctx, user.ID, user.Name)
	res, _ := rows.RowsAffected()
	if res == 0 {
		rp.db.logger.Infoln("Database is empty")
		return 0, nil
	}
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' with his new name '%s'\n", user.ID, user.Name)
	return user.ID, nil
}

func (rp *repo) PartUpdateUserEmail(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if ok, err := rp.IdExist(ctx, user.ID); ok {
		return 0, err
	}

	query := "UPDATE clients SET email = $2 WHERE id = $1 RETURNING id"

	stmt, err := rp.db.sqlDB.PrepareContext(ctx, query)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	defer stmt.Close()

	rows, err := stmt.ExecContext(ctx, user.ID, user.Email)
	res, _ := rows.RowsAffected()
	if res == 0 {
		rp.db.logger.Infoln("Database is empty")
		return 0, nil
	}
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' and new email '%s'\n", user.ID, user.Email)
	return user.ID, nil
}

func (rp *repo) PartUpdateUserPass(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if ok, err := rp.IdExist(ctx, user.ID); ok {
		return 0, err
	}

	if err := user.Valid(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	if err := user.HashingPass(); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	query := "UPDATE clients SET encrypt_password = $2 WHERE id = $1 RETURNING id"
	_, err := rp.db.sqlDB.ExecContext(ctx, query, user.ID, user.Encrypt_Password)
	if err != nil {
		rp.db.logger.Infoln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' with his new password\n", user.ID)
	return user.ID, nil
}

func (rp *repo) IdExist(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := rp.db.sqlDB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM clients WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
