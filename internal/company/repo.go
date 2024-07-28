package company

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"time"
)

type Repository interface {
	CreateCompany(ctx context.Context, name string) (*Company, error)
	GetCompany(ctx context.Context, name string) (*Company, error)
}

type Repo struct {
	model model.Model
}

func (repo *Repo) GetCompany(ctx context.Context, name string) (*Company, error) {
	//TODO implement me
	panic("implement me")
}

func (repo *Repo) CreateCompany(ctx context.Context, name string) (*Company, error) {
	id := util.GenerateUniqueID("CP")
	c := &Company{
		Name: name,
		ID:   id,
	}

	err := repo.model.Transaction(ctx, func(ctx context.Context) error {
		c.CreatedAt = time.Now().UTC()
		err := repo.model.Insert(ctx, "company", c)
		return err
	})
	return c, err
}

func NewRepository(ms *model.Service) *Repo {
	return &Repo{ms}
}
