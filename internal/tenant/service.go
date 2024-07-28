package tenant

import (
	"context"
	"fmt"
)

type TenantUser struct {
	FirstName   string
	LastName    string
	Email       string
	PhoneNumber string
}
type Interface interface {
	CreateTenant(ctx context.Context, tenant *TenantUser, unitID string) (string, error)
}

type Service struct {
	repo Repository
	userService
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTenant(ctx context.Context, tenant *TenantUser, unitID string) (string, error) {
	//TODO: Validate tenant
	id, err := s.repo.CreateTenant(ctx, tenant.FirstName, tenant.LastName, tenant.Email, tenant.PhoneNumber, unitID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return id, nil
}
