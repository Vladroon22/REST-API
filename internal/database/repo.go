package database

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Repo struct {
	db      *DataBase
	timeOut time.Duration
}

func NewRepo(db *DataBase) *Repo {
	return &Repo{
		db:      db,
		timeOut: time.Duration(2) * time.Second,
	}
}

func (rp *Repo) CreateNewUser(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	var id int
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

func (rp *Repo) DeleteUser(c context.Context, id int) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, id); err != nil {
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

func (rp *Repo) UpdateUserFully(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
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

func (rp *Repo) PartUpdateUserName(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
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

func (rp *Repo) PartUpdateUserEmail(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
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

func (rp *Repo) PartUpdateUserPass(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
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

func (rp *Repo) GenerateJWT(c context.Context, email, pass string) (string, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	/*	user, err := rp.GetUser(email, pass)
		if err != nil {
			rp.db.logger.Errorln(err)
			return "", err
		}*/
	user := &User{}
	//	query := "SELECT * FROM clients WHERE email = $1 AND encrypt_password = $2"
	//	err := rp.db.sqlDB.QueryRowContext(ctx, query, email, pass).Scan(&user)
	query := "SELECT id, email, encrypt_password FROM clients WHERE email = $1 AND encrypt_password = $2"
	err := rp.db.sqlDB.QueryRowContext(ctx, query, email, pass).Scan(&user.ID, &user.Email, &user.Encrypt_Password)

	rp.db.logger.Infoln(user.ID)
	rp.db.logger.Infoln(user.Email)
	rp.db.logger.Infoln(user.Encrypt_Password)

	rp.db.logger.Infoln(email)
	rp.db.logger.Infoln(pass)

	if err != nil {
		return "", err
	}

	if err := CmpHashAndPass(user.Encrypt_Password, pass); err != nil {
		rp.db.logger.Errorln(err)
		return "", err
	}

	JWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(), // TTL of token
		IssuedAt:  time.Now().Unix(),
	})

	return JWT.SignedString([]byte("jcuznys^N$74mc8o#9,eijf"))
}

/*
	func (rp *Repo) GetUser(email, pass string) (*User, error) {
		user := &User{}
		query := "SELECT id FROM clients WHERE email = $1 AND encrypt_password = $2"
		err := rp.db.sqlDB.Get(user, query, email, pass)
		if err != nil {
			return nil, err
		}
		rp.db.logger.Infoln(user.Email)
		rp.db.logger.Infoln(user.Encrypt_Password)

		return user, nil
	}
*/
func (rp *Repo) IdExist(ctx context.Context, id int) (int, error) {
	var ID int
	query := "SELECT id FROM clients WHERE id = $1"
	err := rp.db.sqlDB.Get(&ID, query, id)
	if err != nil {
		return 0, err
	}
	return ID, nil
}
