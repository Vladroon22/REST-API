package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	signKey = "jcuznys^N$74mc8o#9,eijf"
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

type MyClaims struct {
	jwt.StandardClaims
	UserId int `json:"id"`
}

func (rp *Repo) CreateNewUser(c context.Context, name, email, pass string) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{
		Name:     name,
		Email:    email,
		Password: pass,
	}
	if err := Valid(user); err != nil {
		return 0, err
	}
	enc_pass, err := HashingPass(pass)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}
	var id int
	query := "INSERT INTO clients (username, email, encrypt_password) VALUES ($1, $2, $3) RETURNING id"
	if err := rp.db.sqlDB.QueryRowContext(ctx, query, user.Name, user.Email, enc_pass).Scan(&id); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infoln("User successfully added")
	return id, nil
}

func (rp *Repo) DeleteUser(c context.Context, id int) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, id); err != nil {
		return 0, err
	}

	query := "DELETE FROM clients WHERE id = $1"
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

	rp.db.logger.Infof("User successfully deleted ID: %d", id)
	return id, nil
}

func (rp *Repo) UpdateUserFully(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
		return 0, err
	}

	enc_pass, err := HashingPass(user.Password)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	query := "UPDATE clients SET username = $1, email = $2, encrypt_password = $3 WHERE id = $4"
	if _, err := rp.db.sqlDB.ExecContext(ctx, query, user.Name, user.Email, enc_pass, user.ID); err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	rp.db.logger.Infof("User successfully updated ID: %d", user.ID)
	return user.ID, nil
}

func (rp *Repo) PartUpdateUserName(c context.Context, user *User) (int, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()

	if _, err := rp.IdExist(ctx, user.ID); err != nil {
		return 0, err
	}

	query := "UPDATE clients SET username = $1 WHERE id = $2"
	if _, err := rp.db.sqlDB.ExecContext(ctx, query, user.Name, user.ID); err != nil {
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

	query := "UPDATE clients SET email = $1 WHERE = $2"
	if _, err := rp.db.sqlDB.ExecContext(ctx, query, user.Email, user.ID); err != nil {
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

	enc_pass, err := HashingPass(user.Password)
	if err != nil {
		rp.db.logger.Errorln(err)
		return 0, err
	}

	query := "UPDATE clients SET encrypt_password = $1 WHERE id = $2"
	if _, err := rp.db.sqlDB.ExecContext(ctx, query, enc_pass, user.ID); err != nil {
		rp.db.logger.Infoln(err)
		return 0, err
	}

	rp.db.logger.Infof("Update the User '%d' with his new password\n", user.ID)
	return user.ID, nil
}

func (rp *Repo) GenerateJWT(c context.Context, email, pass string) (string, error) {
	ctx, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	var id int
	var hash string
	query := "SELECT id, encrypt_password FROM clients WHERE email = $1"
	err := rp.db.sqlDB.QueryRowContext(ctx, query, email).Scan(&id, &hash)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}

	if err := CmpHashAndPass(hash, pass); err != nil {
		rp.db.logger.Errorln(err)
		return "", err
	}

	JWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Hour).Unix(), // TTL of token
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return JWT.SignedString([]byte(signKey))
}

func ValidateToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("Unauthorized")
		}
		return nil, errors.New("Bad Request")
	}

	claims, ok := token.Claims.(*MyClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Unauthorized")
	}

	return claims, nil
}

func (rp *Repo) GetUser(c context.Context, id int) (*User, error) {
	c, cancel := context.WithTimeout(c, rp.timeOut)
	defer cancel()
	user := &User{}
	query := "SELECT id, username, email FROM clients WHERE id = $1"
	err := rp.db.sqlDB.QueryRowContext(c, query, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (rp *Repo) IdExist(ctx context.Context, id int) (int, error) {
	c, cancel := context.WithTimeout(ctx, rp.timeOut)
	defer cancel()
	var ID int
	query := "SELECT id FROM clients WHERE id = $1"
	err := rp.db.sqlDB.QueryRowContext(c, query, id).Scan(&ID)
	if err != nil {
		return 0, errors.New("No such user")
	}
	return ID, nil
}
