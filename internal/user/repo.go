package user

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Repo struct {
	model model.Model
}

func NewRepo(model model.Model) *Repo {
	return &Repo{model: model}
}

func (r *Repo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.model.Transaction(ctx, func(ctx context.Context) error {
		return r.model.Get(ctx,
			&model.SelectQuery{
				TableName: TABLE_NAME,
				Fields:    []string{"*"},
				Where: []model.Condition{
					{Clause: "email", Param: email},
				},
			}, u)
	})
	return u, err
}

func (r *Repo) CreateUser(ctx context.Context, firstName, lastName, email, phoneNumber string) (string, error) {
	userID := ""
	err := r.model.Transaction(ctx, func(ctx context.Context) error {
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
	if err != nil {
		return "", err
	}
	return userID, err
}
