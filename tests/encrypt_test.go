package tests

import (
	"testing"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := &db.User{
		Email:    "user@gmail.com",
		Password: "12345",
	}

	assert.NoError(t, user.HashingPass())
	assert.NotEmpty(t, user.Encrypt_Password)
}
