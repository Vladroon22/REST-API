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
	Password         string `json:"pass"`     // open
}

func CreateNewUser(id int, name, email, password string) *User {
	return &User{
		ID:       1,
		Name:     name,
		Email:    email,
		Password: password,
	}
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
