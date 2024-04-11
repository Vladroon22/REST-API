package database

type User struct {
	ID               int
	Name             string
	Email            string
	Encrypt_Password string // secure
	Password         string // open
}
