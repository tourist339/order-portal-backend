package property

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

type Interface interface {
	CreateProperty(ctx context.Context, address, owner, companyID string) error
	GetProperty(ctx context.Context, propertyID string) (*Property, error)
}

func (s *Service) CreateProperty(ctx context.Context, company, owner, address string) error {
	return s.repo.CreateProperty(ctx, company, owner, address)
}

func (s *Service) GetProperty(ctx context.Context, propertyID string) (*Property, error) {
	return s.repo.GetProperty(ctx, propertyID)
}
