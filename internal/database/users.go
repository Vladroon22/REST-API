package database

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int    `json:"id"`
	Name             string `json:"username"`
	Email            string `json:"email"`
	Encrypt_Password string `json:"enc_pass"` // secure
	Password         string `json:"pass"`
}

func (user *User) HashingPass() error {
	if err := user.Valid(); err != nil {
		return err
	}

	var err error
	if len(user.Password) <= 0 {
		return err
	} else if len(user.Password) > 50 {
		return err
	}

	enc_pass, err := encrypt(user.Password)
	if err != nil {
		return err
	}
	user.Encrypt_Password = enc_pass
	if err := bcrypt.CompareHashAndPassword([]byte(user.Encrypt_Password), []byte(user.Password)); err != nil {
		return err
	}

	return nil
}

func CmpHashAndPass(hash, pass string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		return err
	}
	return nil
}

func encrypt(pass string) (string, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(encrypt), nil
}

func (user *User) Valid() error {
	return validation.ValidateStruct(user, validation.Field(&user.Email, validation.Required, is.Email), validation.Field(&user.Password, validation.Required))
}
