package service

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=cafe_service.go -destination=mocks/cafe_service_mock.go

type CafeService interface {
	GetAllCafes(page, pageSize int) ([]*domain.GetCafeResponse, error)
	GetTotalCafesCount() (int, error)
	GetCafeByID(id primitive.ObjectID) (*domain.GetCafeResponse, error)
	CreateCafe(request *domain.CreateCafeRequest) (*domain.CreateCafeResponse, error)
	UpdateCafe(id primitive.ObjectID, request *domain.UpdateCafeRequest) (*domain.UpdateCafeResponse, error)
	DeleteCafe(id primitive.ObjectID) error
	SearchCafes(query string, page int, pageSize int) ([]*domain.GetCafeResponse, error)
	FilterCafesByTags(tags []string, page int, pageSize int) ([]*domain.GetCafeResponse, error)
}
