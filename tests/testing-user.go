package tests

import (
	"testing"

	db "github.com/Vladroon22/REST-API/internal/database"
)

func CreateUser(t *testing.T) *db.User {
	return &db.User{
		Name:     "nickname",
		Email:    "example@gmail.com",
		Password: "12345678",
	}
}
