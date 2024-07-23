package database

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"pass"`
}

func HashingPass(password string) (string, error) {
	enc_pass, err := encrypt(password)
	if err != nil {
		return "", err
	}
	return enc_pass, nil
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

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(email)
}

func Valid(user *User) error {
	if user.Password == "" {
		return errors.New("password cant't be blank")
	} else if len(user.Password) >= 50 {
		return errors.New("password cant't be more than 50 symbols")
	} else if user.Name == "" {
		return errors.New("username cant't be blank")
	} else if user.Email == "" {
		return errors.New("email can't be blank")
	}

	if ok := validateEmail(user.Email); !ok {
		return errors.New("wrong email input")
	}
	return nil
}
