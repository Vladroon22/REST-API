package database

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int
	Name             string
	Email            string
	Encrypt_Password string // secure
	Password         string // open
}

func (user *User) HashingPass() error {
	if len(user.Password) > 0 {
		enc_pass, err := encrypt(user.Password)
		if err != nil {
			return err
		}
		user.Encrypt_Password = enc_pass
	}

	return nil
}

func encrypt(pass string) (string, error) {
	encrypt, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(encrypt), nil
}
