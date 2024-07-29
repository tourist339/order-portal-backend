package role

import (
	"backend/internal/user"
	"context"
	"fmt"
)

type Interface interface {
	CreateWorker(ctx context.Context, workerRoleType string, roleID string, propertyID string) error
	CreateTenant(ctx context.Context, tenant *user.BasicUser, propertyID string) (string, error)
}

type Service struct {
	repo        Repository
	userService user.Interface
}

func NewService(repo Repository, userService user.Interface) *Service {
	return &Service{repo: repo, userService: userService}
}
func (s *Service) CreateWorker(ctx context.Context, workerRoleType string, roleID string, propertyID string) error {
	_, err := s.repo.CreateRole(ctx, "worker:"+workerRoleType, roleID, propertyID)
	return err
}

func (s *Service) CreateTenant(ctx context.Context, tenant *user.BasicUser, propertyID string) (string, error) {
	//TODO: Validate tenant
	userID, err := s.userService.CreateOrGetUser(ctx, tenant.FirstName, tenant.LastName, tenant.Email, tenant.PhoneNumber)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	roleID, err := s.repo.CreateRole(ctx, ROLE_TYPE_TENANT, userID, propertyID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return roleID, nil
}
