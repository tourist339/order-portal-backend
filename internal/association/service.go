package association

import "context"

type Interface interface {
	CreateWorkerAssociation(ctx context.Context, workerRoleType string, roleID string, propertyID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWorkerAssociation(ctx context.Context, workerRoleType string, roleID string, propertyID string) error {
	return s.repo.CreateAssociation(ctx, "worker:"+workerRoleType, roleID, propertyID)
}
