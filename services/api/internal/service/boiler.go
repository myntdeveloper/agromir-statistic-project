package service

import (
	"context"
	"errors"

	"github.com/myntdeveloper/agromir_boiler_server/internal/model"
	"github.com/myntdeveloper/agromir_boiler_server/internal/repo"
)

var (
	ErrWarehouseNotFound = errors.New("warehouse not found")
	ErrBoilerNotFound    = errors.New("boiler not found")
	ErrPlantNotFound     = errors.New("plant not found")
)

type BoilerService struct {
	repo *repo.BoilerRepo
}

func NewBoilerService(repo *repo.BoilerRepo) *BoilerService {
	return &BoilerService{repo: repo}
}

func (s *BoilerService) ListWarehouses(ctx context.Context) ([]model.WarehouseSummary, error) {
	return s.repo.ListWarehouses(ctx)
}

func (s *BoilerService) ListBoilers(ctx context.Context, warehouseID int) ([]model.BoilerListItem, error) {
	exists, err := s.repo.WarehouseExists(ctx, warehouseID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrWarehouseNotFound
	}
	return s.repo.ListBoilersByWarehouse(ctx, warehouseID)
}

func (s *BoilerService) GetBoiler(ctx context.Context, id int) (*model.BoilerDetail, error) {
	boiler, err := s.repo.GetBoiler(ctx, id)
	if err != nil {
		return nil, err
	}
	if boiler == nil {
		return nil, ErrBoilerNotFound
	}
	return boiler, nil
}

func (s *BoilerService) GetTrend(ctx context.Context, boilerID int) (*model.BoilerTrend, error) {
	exists, err := s.repo.BoilerExists(ctx, boilerID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrBoilerNotFound
	}
	points, err := s.repo.GetTrend(ctx, boilerID)
	if err != nil {
		return nil, err
	}
	return &model.BoilerTrend{BoilerID: boilerID, Points: points}, nil
}
