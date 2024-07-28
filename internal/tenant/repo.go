package tenant

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"fmt"
	"time"
)

type Repository interface {
	CreateTenant(ctx context.Context, userID, unitID string) (string, error)
}

type Repo struct {
	model model.Model
}

func NewRepo(model model.Model) *Repo {
	return &Repo{model: model}
}
func (r *Repo) CreateTenant(ctx context.Context, userID, unitID string) (string, error) {
	id := util.GenerateUniqueID("TN")
	return id, r.model.Transaction(ctx, func(ctx context.Context) error {
		tenant := &Tenant{
			ID:        id,
			UserID:    userID,
			UnitID:    unitID,
			CreatedAt: time.Now().UTC(),
		}
		err := r.model.Insert(ctx, TABLE_NAME, tenant)
		if err != nil {
			fmt.Println(fmt.Errorf("Error inserting tenant: %v", err.Error()))
		}
		return err
	})
}
