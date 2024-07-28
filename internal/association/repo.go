package association

import (
	"backend/internal/model"
	"backend/internal/util"
	"context"
	"time"
)

type Repository interface {
	CreateAssociation(ctx context.Context, roleType string, roleID string, propertyID string) error
}

type Repo struct {
	model model.Model
}

func NewRepo(model model.Model) *Repo {
	return &Repo{model: model}
}

func (r *Repo) CreateAssociation(ctx context.Context, roleType string, roleID string, propertyID string) error {
	return r.model.Transaction(ctx, func(ctx context.Context) error {
		association := &UserPropertyAssociation{
			ID:         util.GenerateUniqueID("UPA"),
			RoleID:     roleID,
			PropertyID: propertyID,
			RoleType:   roleType,
			CreatedAt:  time.Now().UTC(),
		}
		err := r.model.Insert(ctx, TABLE_NAME, association)
		return err
	})
}
