package tests

import (
	"testing"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := db.CreateUserForTest(t)

	assert.NoError(t, user.HashingPass())
	assert.NotEmpty(t, user.Encrypt_Password)
	assert.NoError(t, user.Valid())
}

func TestDataValid(t *testing.T) {
	Cases := []struct {
		username string
		valid    bool
		us       func() *db.User
	}{
		{
			username: "vlad",
			us:       func() *db.User { return db.CreateUserForTest(t) }, // !!!
			valid:    true,
		},
		{
			username: "null email",
			us: func() *db.User {
				user := db.CreateUserForTest(t) /// !!!
				user.Email = " "
				return user

			},
			valid: false,
		},
		{
			username: "null password",
			us: func() *db.User {
				user := db.CreateUserForTest(t) /// !!!
				user.Password = ""
				return user
			},
			valid: false,
		},
	}

	for _, cases := range Cases {
		t.Run(cases.username, func(t *testing.T) {
			if cases.valid {
				assert.NoError(t, cases.us().Valid())
			} else {
				assert.Error(t, cases.us().Valid())
			}
		})
	}
}
