package roles

import (
	"backend/internal/user"
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
	repo        Repository
	userService user.Interface
}

func NewService(repo Repository, userService user.Interface) *Service {
	return &Service{repo: repo, userService: userService}
}

func (s *Service) CreateTenant(ctx context.Context, tenant *TenantUser, unitID string) (string, error) {
	//TODO: Validate tenant
	userID, err := s.userService.CreateUser(ctx, tenant.FirstName, tenant.LastName, tenant.Email, tenant.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	id, err := s.repo.CreateTenant(ctx, userID, unitID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return id, nil
}
