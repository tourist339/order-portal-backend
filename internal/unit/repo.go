package unit

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"fmt"
	"time"
)

type Repository interface {
	CreateUnit(ctx context.Context, propertyID string, unitName string) (string, error)
	GetUnit(ctx context.Context, propertyID string, unitName string) (*Unit, error)
	AddTenantToUnit(ctx context.Context, tenantID []string, unitID string) error
}

type Repo struct {
	store model.Model
}

func NewRepo(store model.Model) *Repo {
	return &Repo{store: store}
}

func (r *Repo) CreateUnit(ctx context.Context, propertyID string, unitName string) (string, error) {
	id := util.GenerateUniqueID("UT")
	err := r.store.Transaction(ctx, func(ctx context.Context) error {
		u := &Unit{
			ID:         id,
			Name:       unitName,
			PropertyID: propertyID,
			CreatedAt:  time.Now().UTC(),
		}
		err := r.store.Insert(ctx, TABLE_NAME, u)
		return err
	})

	return id, err
}

func (r *Repo) GetUnit(ctx context.Context, propertyID string, unitName string) (*Unit, error) {
	u := &Unit{}
	err := r.store.Transaction(ctx, func(ctx context.Context) error {
		return r.store.Get(ctx, &model.SelectQuery{
			TableName: TABLE_NAME,
			Fields:    []string{"*"},
			Where: []model.Condition{
				{
					Clause: "property_id",
					Param:  propertyID,
				},
				{
					Clause: "name",
					Param:  unitName,
				},
			},
		}, u)
	})
	if err != nil {
		return nil, err
	}
	return u, err
}
func (r *Repo) AddTenantToUnit(ctx context.Context, tenantIDs []string, unitID string) error {
	err := r.store.Transaction(ctx, func(ctx context.Context) error {
		u := &Unit{}
		err := r.store.GetByID(ctx, unitID, TABLE_NAME, []string{"*"}, u)
		if err != nil {
			fmt.Println("Error getting unit", err)
			return err
		}
		fmt.Println(fmt.Sprintf("Unit fetched %v", u))
		u.TenantIDs = tenantIDs
		err = r.store.Update(ctx, TABLE_NAME, u)
		return err
	})
	return err
}
