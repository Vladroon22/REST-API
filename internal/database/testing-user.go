package database

import "testing"

func CreateUser(t *testing.T) *User {
	return &User{
		Email:    "example@gmail.com",
		Password: "12345678",
	}
}
