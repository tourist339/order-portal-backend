package property

import (
	"backend/internal/company"
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"time"
)

const TABLE_NAME = "property"

type Repo struct {
	model    model.Model
	compRepo company.Repository
}
type Repository interface {
	CreateProperty(ctx context.Context, company, address string) error
	GetProperty(ctx context.Context, propertyId string) (*Property, error)
}

func NewRepository(m model.Model, compRepo company.Repository) *Repo {
	return &Repo{m, compRepo}
}

type Company struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func (r *Repo) CreateProperty(ctx context.Context, address, owner, companyID string) error {

	err := r.model.Transaction(ctx, func(ctx context.Context) error {
		prop := &Property{
			ID:        util.GenerateUniqueID("PR"),
			CompanyID: companyID,
			Address:   address,
			//Owner:     sql.NullString{String: owner},
			//TeamIDs:   []string{"1abc", "232"},
			CreatedAt: time.Now().UTC(),
		}
		_, err := r.model.Insert(ctx, "property", prop)
		if err != nil {
			println(err)
		}
		return err
	})
	return err

}

func (r *Repo) GetProperty(ctx context.Context, id string) (*Property, error) {
	p := &Property{}
	err := r.model.Transaction(ctx, func(ctx context.Context) error {
		return r.model.GetByID(ctx, id, TABLE_NAME, []string{"*"}, p)
	})
	if err != nil {
		return nil, err
	}
	return p, nil

}
