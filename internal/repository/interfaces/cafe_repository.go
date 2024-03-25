package repository

import (
	"vendors/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source=cafe_repository.go -destination=mocks/cafe_repository_mock.go

type CafeRepository interface {
	GetAllCafes(page, pageSize int) ([]*domain.GetCafeResponse, error)
	GetTotalCafesCount() (int, error)
	GetCafeByID(id primitive.ObjectID) (*domain.GetCafeResponse, error)
	CreateCafe(request *domain.CreateCinemaRequest) (*domain.CreateCafeResponse, error)
	UpdateCafe(id primitive.ObjectID, request *domain.UpdateCinemaRequest) (*domain.UpdateCafeResponse, error)
	DeleteCafe(id primitive.ObjectID) error
	SearchCafes(query string, page int, pageSize int) ([]*domain.GetCafeResponse, error)
	FilterCafesByTags(tags []string, page int, pageSize int) ([]*domain.GetCafeResponse, error)
}
