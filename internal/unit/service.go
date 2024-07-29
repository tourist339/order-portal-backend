package unit

import (
	"backend/internal/model"
	"backend/internal/property"
	"backend/internal/role"
	"context"
	"errors"
	"fmt"
)

var ErrUnitAlreadyExists = errors.New("unit already exists")

type Interface interface {
	CreateUnit(ctx context.Context, unitIdentifier string, propertyID string) (string, error)
	GetUnit(ctx context.Context, unitIdentifier string, propertyID string) (*Unit, error)
	GetUnitsByProperty(ctx context.Context, propertID string) ([]*Unit, error)
}

type Service struct {
	repo            Repository
	roleService     role.Interface
	propertyService property.Interface
}

func NewService(repo Repository, roleService role.Interface, propertyService property.Interface) *Service {
	return &Service{repo: repo, roleService: roleService, propertyService: propertyService}
}

func (s *Service) CreateUnit(ctx context.Context, propertyID, unitIdentifier string) (string, error) {
	//Check Property Exists
	_, err := s.propertyService.GetProperty(ctx, propertyID)
	if err != nil {
		fmt.Println("Error getting Property by", propertyID, err)
		return "", err
	}
	u, err := s.repo.GetUnit(ctx, propertyID, unitIdentifier)
	if err == nil {
		return u.ID, nil
	}
	if !errors.Is(err, model.ErrNotFound) {
		fmt.Println("Error getting unit", err)
		return "", err
	}
	//Create Unit
	unitID, err := s.repo.CreateUnit(ctx, propertyID, unitIdentifier)
	if err != nil {
		fmt.Println("Error creating unit", err)
		return "", err
	}
	//TODO: Create Tenant in goroutines

	return unitID, nil
}

func (s *Service) GetUnit(ctx context.Context, propertyID, unitIdentifier string) (*Unit, error) {
	u, err := s.repo.GetUnit(ctx, propertyID, unitIdentifier)
	if err != nil {
		fmt.Println("Error getting unit", err)
		return nil, err
	}
	return u, nil
}

func (s *Service) GetUnitsByProperty(ctx context.Context, propertyID string) ([]*Unit, error) {
	//units, err := s.repo.GetUnitsByProperty(ctx, propertyID)
	//if err != nil {
	//	fmt.Println("Error getting units by property", err)
	//	return nil, err
	//}
	return nil, nil
}
