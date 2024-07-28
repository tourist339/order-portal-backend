package user

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error)
}

type Repo struct {
	model model.Model
}

func NewRepo(model model.Model) *Repo {
	return &Repo{model: model}
}

func (r *Repo) CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error) {
	userID := ""
	r.model.Transaction(ctx, func(ctx context.Context) error {
		userID = util.GenerateUniqueID("U")
		u := &User{
			ID:        userID,
			FirstName: firstName,
			LastName:  lastName,
			Email:     email,
			Phone:     phoneNumber,
		}
		return r.model.Insert(ctx, TABLE_NAME, u)
	})
}
