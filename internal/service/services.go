package service

import (
	"context"

	db "github.com/Vladroon22/REST-API/internal/database"
)

type Accounts interface {
	CreateNewUser(c context.Context, user *db.User) (int, error)
	DeleteUser(c context.Context, id int) (int, error)
	UpdateUserFully(c context.Context, user *db.User) (int, error)
	PartUpdateUserName(c context.Context, user *db.User) (int, error)
	PartUpdateUserEmail(c context.Context, user *db.User) (int, error)
	PartUpdateUserPass(c context.Context, user *db.User) (int, error)
	GenerateJWT(email, pass string) (string, error)
	GetUser(email, pass string) (*db.User, error)
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
