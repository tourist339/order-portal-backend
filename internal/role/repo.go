package role

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"time"
)

type Repository interface {
	CreateRole(ctx context.Context, roleType string, userID string, propertyID string) (string, error)
}

type Repo struct {
	model model.Model
}

func NewRepo(model model.Model) *Repo {
	return &Repo{model: model}
}

func (r *Repo) CreateRole(ctx context.Context, roleType string, userID string, propertyID string) (string, error) {
	roleID := util.GenerateUniqueID("RL")

	return roleID, r.model.Transaction(ctx, func(ctx context.Context) error {
		association := &Role{
			ID:         roleID,
			UserID:     userID,
			PropertyID: propertyID,
			RoleType:   roleType,
			CreatedAt:  time.Now().UTC(),
		}
		err := r.model.Insert(ctx, TABLE_NAME, association)
		return err
	})
}
