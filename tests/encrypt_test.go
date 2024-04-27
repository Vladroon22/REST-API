package tests

import (
	"testing"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := db.CreateUser()

	assert.NoError(t, user.HashingPass())
	assert.NotEmpty(t, user.Encrypt_Password)
	assert.NoError(t, user.Valid())
}

func TestDataValid(t *testing.T) {
	Cases := []struct {
		username string
		valid    bool
		us       func(*db.User)
		us1      func() *db.User
	}{
		{
			username: "vlad",
			valid:    true,
			us:       func(*db.User) { TestUser(t) },
		},
		{
			username: "null email",
			valid:    false,
			us1: func() *db.User {
				user := db.CreateUser()
				user.Email = ""
				return user
			},
		},
		{
			username: "null password",
			valid:    false,
			us1: func() *db.User {
				user := db.CreateUser()
				user.Password = ""
				return user
			},
		},
	}

	for _, cases := range Cases {
		t.Run(cases.username, func(t *testing.T) {
			user := db.CreateUser()
			cases.us(user)
			if cases.valid {
				assert.NoError(t, cases.us1().Valid())
				assert.NoError(t, user.Valid())
			} else {
				assert.Error(t, cases.us1().Valid())
				assert.Error(t, user.Valid())
			}
		})
	}
}
