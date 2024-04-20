package tests

import (
	"testing"

	"github.com/Vladroon22/REST-API/config"
	"github.com/Vladroon22/REST-API/internal/database"
	"github.com/stretchr/testify/assert"
)

func setupDB() (*database.DataBase, error) {
	conf := config.CreateConfig()
	conf.Host = "localhost"
	conf.Username = "postgres"
	conf.Password = ""
	conf.DBname = "db_test"
	conf.SSLmode = "sslmode=disable"
	conf.Port = "5432"
	conf.Addr_PORT = ""
	db := database.NewDB(conf)
	if err := db.ConfigDB(); err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateOfUser(t *testing.T) {
	t.Helper()
	db, err := setupDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.CloseDB()

	user, err := db.SetUser().CreateNewUser(&database.User{
		Name:  "DVD",
		Email: "dvdv123@gmail.com",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)
}
