package database

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
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

func (user *User) Valid() error {
	return validation.ValidateStruct(user, validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(8, 50)))
}

func CreateUserForTest(t *testing.T) *User {
	return &User{
		Email:    "example@gmail.com",
		Password: "12345678",
	}
}
