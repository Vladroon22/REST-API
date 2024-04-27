package tests

import (
	"testing"

	db "github.com/Vladroon22/REST-API/internal/database"
	"github.com/stretchr/testify/assert"
)

func TestDataValid(t *testing.T) {
	Cases := []struct {
		username string
		valid    bool
		us       func() *db.User
	}{
		{
			username: "vlad",
			valid:    true,
			us:       func() *db.User { return db.CreateUser(t) },
		},
		{
			username: "null email",
			valid:    false,
			us: func() *db.User {
				user := db.CreateUser(t)
				user.Email = ""
				return user
			},
		},
		{
			username: "null password",
			valid:    false,
			us: func() *db.User {
				user := db.CreateUser(t)
				user.Password = ""
				return user
			},
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
