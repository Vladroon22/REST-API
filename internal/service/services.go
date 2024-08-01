package service

import (
	"context"

	db "github.com/Vladroon22/REST-API/internal/database"
)

type Accounts interface {
	CreateNewUser(c context.Context, name, email, pass string) (int, error)
	DeleteUser(c context.Context, id int) (int, error)
	UpdateUserFully(c context.Context, user *db.User) (int, error)
	PartUpdateUserName(c context.Context, user *db.User) (int, error)
	PartUpdateUserEmail(c context.Context, user *db.User) (int, error)
	PartUpdateUserPass(c context.Context, user *db.User) (int, error)
	GenerateJWT(c context.Context, email, pass string) (string, error)
	GetUser(c context.Context, id int) (*db.User, error)
	IdExist(ctx context.Context, id int) (int, error)
}

type Service struct {
	Accounts
}

func NewService(repos *db.Repo) *Service {
	return &Service{
		Accounts: repos,
	}
}
